package main

import (
	"Huffman-Compression-Go/lib"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func DecodeText(enc_file string) {
	data, err := os.ReadFile(enc_file)
	if err != nil {
		panic(err)
	}
	topo_length := int(binary.LittleEndian.Uint32(data[0:4]))
	total_chars := int(binary.LittleEndian.Uint32(data[4:8]))

	tree_topology := data[8 : topo_length+8]
	tree := lib.DecodeTopology(tree_topology)

	root := tree
	databs := lib.BitSet(data[topo_length+8:])

	decoded := make([]rune, total_chars)
	ptr := 0
	for idx := 0; idx < 8*databs.Length() && ptr < total_chars; idx++ {
		if databs.GetBit(idx) {
			root = root.Right
		} else {
			root = root.Left
		}
		if root.IsLeaf {
			decoded[ptr] = root.Token
			ptr++
			root = tree
		}
	}
	out_f := enc_file[:strings.IndexByte(enc_file, '.')] + "_extracted.txt"
	err = os.WriteFile(out_f, []byte(string(decoded)), 0666)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nUncompressed file into", out_f)
}

func main() {
	args := os.Args[1:]
	filename := args[0]

	DecodeText(filename)
}
