package main

import (
	"Huffman-Compression-Go/lib"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// Write compressed data into file
func EncodeText(filename string) {

	bdata, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	text := string(bdata)

	frequency_table := make(map[rune]int)
	for _, ch := range text {
		frequency_table[ch]++
	}
	// TEST
	// for tok, cnt := range frequency_table {
	// 	fmt.Println("Token :", tok, "Count :", cnt)
	// }
	// TEST

	tree := lib.BuildHuffmanTree(frequency_table)

	enc_table := make(map[rune]string)
	lib.GetTokenCodes("", tree, enc_table)

	// Encode text into a byte array
	size := 0
	for tok, freq := range frequency_table {
		size += freq * len(enc_table[tok])
	}
	buffer := lib.NewBitSet(size)
	buff_ptr := 0

	for _, token := range text {
		code := enc_table[token]
		for _, bit := range code {
			if bit == '0' {
				buffer.SetBit(buff_ptr, false)
			} else {
				buffer.SetBit(buff_ptr, true)
			}
			buff_ptr++
		}
	}

	// We write a header section to our compressed file in order to decode it. Taken from https://engineering.purdue.edu/ece264/17au/hw/HW13?alt=huffman, this includes :
	// Total no of characters storing tree topology, Total no of characters in the original file
	// Tree topology
	// followed by the body : byte array
	fname := filename[:strings.IndexByte(filename, '.')]
	f, err := os.Create(fname + ".enc")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Generate header for file, ie the postfix topology of tree
	var headerbuff []byte
	lib.GenerateTopology(tree, &headerbuff)
	headerbuff = append(headerbuff, 0)

	total_chars := uint32(len(text))
	topo_length := uint32(len(headerbuff))

	wordbuff := make([]byte, 8)
	binary.LittleEndian.PutUint32(wordbuff, topo_length)
	binary.LittleEndian.PutUint32(wordbuff[4:], total_chars)

	_, err = f.Write(wordbuff)
	if err != nil {
		fmt.Println("Writing header bits..")
		panic(err)
	}
	_, err = f.Write(headerbuff)
	if err != nil {
		fmt.Println("Writing tree topology..")
		panic(err)
	}

	_, err = f.Write(buffer)
	if err != nil {
		fmt.Println("Writing encoded data...")
		panic(err)
	}

	fmt.Println("\nWrote compressed text into " + fname + ".enc !")
}

func main() {
	args := os.Args[1:]
	filename := args[0]

	EncodeText(filename)
}
