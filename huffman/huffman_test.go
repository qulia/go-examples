package huffman

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestHuffmanEncoding(t *testing.T) {
	testRoundTrip(t, "aab")
	testRoundTrip(t, "go go gophers")
	testRoundTrip(t, "Huffman coding is a data compression algorithm.")
}

func testRoundTrip(t *testing.T, input string) {
	res, hTree := Encode(input)
	fmt.Printf("Hufman code: %s", res)
	assert.Equal(t, input, Decode(res, hTree))
}
