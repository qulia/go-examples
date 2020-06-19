package huffman

import (
	"fmt"
	"strings"

	"github.com/qulia/go-log/log"

	"github.com/qulia/go-qulia/lib"
	"github.com/qulia/go-qulia/lib/heap"
	"github.com/qulia/go-qulia/lib/tree"
)

type Item struct {
	Key  string
	Freq int
}

// Golang implementation of huffman algorithm: https://en.wikipedia.org/wiki/Huffman_coding
// Note that the values are still using string in the real implementation 0 and 1s would be packed as bits for
// compression
func Encode(input string) (string, *tree.Node) {
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

func Decode(encodedString string, hTree *tree.Node) string {
	input := []rune(encodedString)
	builder := strings.Builder{}
	decode(input, 0, &builder, hTree, hTree)

	return builder.String()
}

func decode(input []rune, index int, builder *strings.Builder, root *tree.Node, hTree *tree.Node) {
	if root.Left == nil && root.Right == nil {
		builder.WriteString(root.Data.(Item).Key)
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

func getHuffmanEncodingMap(root *tree.Node) map[string]string {
	eMap := make(map[string]string)
	generateMap(root, "", eMap)
	log.V("Emap: %v", eMap)
	return eMap
}

func generateMap(root *tree.Node, current string, eMap map[string]string) {
	if root.Left == nil && root.Right == nil {
		eMap[root.Data.(Item).Key] = current
		return
	}

	generateMap(root.Left, current+"0", eMap)
	generateMap(root.Right, current+"1", eMap)
}

func printTree(root *tree.Node) {
	var result []string
	tree.VisitInOrder(root, func(elem interface{}) {
		var elemString string
		if elem == nil {
			elemString = "nil"
		} else {
			elemString = fmt.Sprintf("%v", elem.(Item))
		}

		result = append(result, elemString)
	})

	log.V("Huffman tree: %v", result)
}

func getHuffmanTree(input string) *tree.Node {
	// Create freq map
	inputR := []rune(input)
	freq := make(map[rune]int)
	for i := 0; i < len(input); i++ {
		if _, ok := freq[inputR[i]]; !ok {
			freq[inputR[i]] = 0
		}

		freq[inputR[i]]++
	}

	hHeap := heap.NewMinHeap(nil, func(first, second interface{}) int {
		// Define comparison function
		firstData := first.(*tree.Node).Data.(Item)
		secondData := second.(*tree.Node).Data.(Item)

		return lib.IntCompFunc(firstData.Freq, secondData.Freq)
	})

	// Build heap from the map
	for k, v := range freq {
		node := tree.NewNode(Item{
			Key:  string(k),
			Freq: v,
		})

		hHeap.Insert(node)
	}

	// Until there is one element left pick min 2, combine and push new tree
	for hHeap.Size() > 1 {
		one := hHeap.Extract()
		two := hHeap.Extract()

		itemOne := one.(*tree.Node).Data.(Item)
		itemTwo := two.(*tree.Node).Data.(Item)
		node := tree.NewNode(Item{
			Key:  itemOne.Key + itemTwo.Key,
			Freq: itemOne.Freq + itemTwo.Freq,
		})
		node.Left = one.(*tree.Node)
		node.Right = two.(*tree.Node)

		hHeap.Insert(node)
	}

	return hHeap.Extract().(*tree.Node)
}
