package lib

import (
	"container/heap"
	"unicode/utf8"
	"unsafe"
)

// Huffman node struct and Priority Queue impl
type HuffmanNode struct {
	Weight      int
	IsLeaf      bool
	Token       rune
	Left, Right *HuffmanNode
}

// Min-Heap of HuffmanNode pointers impl
type NodeHeap []*HuffmanNode

func (h NodeHeap) Len() int {
	return len(h)
}
func (h NodeHeap) Less(i, j int) bool {
	return h[i].Weight < h[j].Weight
}
func (h NodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *NodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*HuffmanNode))
}
func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x

	/* Experimental : Doesnt work | Tried to generate a frequency table that reproduces the same Huffman tree no matter how decoded
	old := *h
	n := len(old)
	last := old[n - 1]
	ptr  := n - 1

	for i := n - 2; i >= 0; i -- {
	  if old[i].Weight != last.Weight {
	    break
	  }
	  if old[i].Token < last.Token {
	    ptr = i
	    last = old[i]
	  }
	}
	*h = append(old[:ptr], old[ptr + 1:]...)
	return last
	*/
}

// Stack of HuffmanNode pointers impl
type NodeStack []*HuffmanNode

func (s *NodeStack) Len() int {
	return len(*s)
}
func (s *NodeStack) Push(node *HuffmanNode) {
	*s = append(*s, node)
}
func (s *NodeStack) Pop() *HuffmanNode {
	if s.Len() == 0 {
		return nil
	} else {
		idx := len(*s) - 1
		element := (*s)[idx]
		*s = (*s)[:idx]
		return element
	}
}

// Bitset data structure for writing and reading bits
type BitSet []uint8

func NewBitSet(l int) BitSet {
	return make(BitSet, (l+7)/8)
}
func (b BitSet) Length() int {
	return len(b)
}
func (b BitSet) GetBit(idx int) bool {
	bytepos, bitpos := idx/8, uint8(idx%8)
	return (b[bytepos] & (uint8(1) << bitpos)) != 0
}
func (b BitSet) SetBit(idx int, val bool) {
	bytepos, bitpos := idx/8, uint8(idx%8)
	if val {
		b[bytepos] |= (uint8(1) << bitpos)
	} else {
		b[bytepos] &= ^(uint8(1) << bitpos)
	}
}

// Building the Huffman Tree using Frequency table
func BuildHuffmanTree(table map[rune]int) *HuffmanNode {
	// Make Heap out of frequency table : Create list of leaf HuffmanNodes and then for each node heapify
	nodeheap := &NodeHeap{}
	heap.Init(nodeheap)
	for tok, freq := range table {
		node := &HuffmanNode{freq, true, tok, nil, nil}
		heap.Push(nodeheap, node)
	}

	// Build tree
	for nodeheap.Len() != 1 {
		a := heap.Pop(nodeheap).(*HuffmanNode)
		b := heap.Pop(nodeheap).(*HuffmanNode)
		merged := &HuffmanNode{a.Weight + b.Weight, false, rune(2147483647), a, b}
		heap.Push(nodeheap, merged)
	}
	return heap.Pop(nodeheap).(*HuffmanNode)
}

// Get binary codes for characters/tokens in text using Huffman tree
func GetTokenCodes(code string, root *HuffmanNode, table map[rune]string) {
	if root.IsLeaf == true {
		table[root.Token] = code
		return
	}
	GetTokenCodes(code+"0", root.Left, table)
	GetTokenCodes(code+"1", root.Right, table)
}

// Auxiliary Rune to Byte array function taken from https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func RuneToBytes(num rune) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

// This function encodes the tree into a Writable byte buffer. Have modified code so that all runes are encoded as 4-byte arrays.
func GenerateTopology(root *HuffmanNode, buffer *[]byte) {
	if root.IsLeaf {
		*buffer = append(*buffer, 1)
		b := RuneToBytes(root.Token)
		for i := 0; i < 4; i++ {
			*buffer = append(*buffer, b[i])
		}
		return
	}
	GenerateTopology(root.Left, buffer)
	GenerateTopology(root.Right, buffer)
	*buffer = append(*buffer, 0)
}

// This function decodes the tree from the topology byte buffer. Runes are decoded as 4 byte arrays
func DecodeTopology(header []byte) *HuffmanNode {
	stack := &NodeStack{}
	ptr := 0
	for true {
		if header[ptr] == 1 {
			r, _ := utf8.DecodeRune(header[ptr+1 : ptr+5])
			node := &HuffmanNode{0, true, r, nil, nil}
			stack.Push(node)
			ptr += 5
		} else if header[ptr] == 0 {
			if stack.Len() > 1 {
				r := stack.Pop()
				l := stack.Pop()
				parent := &HuffmanNode{0, false, rune(2147483647), l, r}
				stack.Push(parent)
				ptr += 1
			} else {
				break
			}
		}
	}
	return stack.Pop()
}
