package huffman

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qulia/go-log/log"

	"github.com/qulia/go-qulia/lib/heap"
	"github.com/qulia/go-qulia/lib/tree"
)

type KeyMap struct {
	Key  string
	Freq int
}

type HeapItem struct {
	n *tree.Node[string]
}

func (hi HeapItem) Compare(other HeapItem) int {
	freqhi := (&KeyMap{}).FromString(hi.n.Data).Freq
	freqother := (&KeyMap{}).FromString(other.n.Data).Freq
	if freqhi < freqother {
		return -1
	} else if freqhi > freqother {
		return 1
	}
	return 0
}

// qo-qulia/lib/tree does not have flex implementation yet so
// as workaround convert the values to string to store in the nodes
func (km KeyMap) ToString() string {
	return fmt.Sprintf("%s:%d", km.Key, km.Freq)
}

func (km *KeyMap) FromString(s string) *KeyMap {
	parts := strings.Split(s, ":")
	km.Key = parts[0]
	km.Freq, _ = strconv.Atoi(parts[1])
	return km
}

// Golang implementation of huffman algorithm: https://en.wikipedia.org/wiki/Huffman_coding
// Note that the values are still using string in the real implementation 0 and 1s would be packed as bits for
// compression
func Encode(input string) (string, *tree.Node[string]) {
	// Create huffman tree and map
	hTree := getHuffmanTree(input)
	printTree(hTree)
	hMap := getHuffmanEncodingMap(hTree)
	builder := strings.Builder{}
	for i := 0; i < len(input); i++ {
		builder.WriteString(hMap[string(input[i])])
	}
	return builder.String(), hTree
}

func Decode(encodedString string, hTree *tree.Node[string]) string {
	input := []rune(encodedString)
	builder := strings.Builder{}
	decode(input, 0, &builder, hTree, hTree)

	return builder.String()
}

func decode(input []rune, index int, builder *strings.Builder, root *tree.Node[string], hTree *tree.Node[string]) {
	if root.Left == nil && root.Right == nil {
		builder.WriteString((&KeyMap{}).FromString(root.Data).Key)
		decode(input, index, builder, hTree, hTree)
		return
	}

	if index == len(input) {
		return
	}
	if input[index] == '0' {
		decode(input, index+1, builder, root.Left, hTree)
	} else {
		decode(input, index+1, builder, root.Right, hTree)
	}
}

func getHuffmanEncodingMap(root *tree.Node[string]) map[string]string {
	eMap := make(map[string]string)
	generateMap(root, "", eMap)
	log.V("Emap: %v", eMap)
	return eMap
}

func generateMap(root *tree.Node[string], current string, eMap map[string]string) {
	if root.Left == nil && root.Right == nil {
		eMap[(&KeyMap{}).FromString(root.Data).Key] = current
		return
	}

	generateMap(root.Left, current+"0", eMap)
	generateMap(root.Right, current+"1", eMap)
}

func printTree(root *tree.Node[string]) {
	var result []string
	tree.VisitInOrder(root, func(elem string) {
		result = append(result, elem)
	})

	log.V("Huffman tree: %v", result)
}

func getHuffmanTree(input string) *tree.Node[string] {
	// Create freq map
	inputR := []rune(input)
	freq := make(map[rune]int)
	for i := 0; i < len(input); i++ {
		if _, ok := freq[inputR[i]]; !ok {
			freq[inputR[i]] = 0
		}

		freq[inputR[i]]++
	}

	hHeap := heap.NewMinHeapFlex[HeapItem](nil)
	// Build heap from the map
	for k, v := range freq {
		node := tree.NewNode(KeyMap{
			Key:  string(k),
			Freq: v,
		}.ToString())

		hHeap.Insert(HeapItem{n: node})
	}

	// Until there is one element left pick min 2, combine and push new tree
	for hHeap.Size() > 1 {
		one := hHeap.Extract()
		two := hHeap.Extract()

		itemOne := (&KeyMap{}).FromString(one.n.Data)
		itemTwo := (&KeyMap{}).FromString(two.n.Data)
		node := tree.NewNode(KeyMap{
			Key:  itemOne.Key + itemTwo.Key,
			Freq: itemOne.Freq + itemTwo.Freq,
		}.ToString())
		node.Left = one.n
		node.Right = two.n

		hHeap.Insert(HeapItem{n: node})
	}

	return hHeap.Extract().n
}
