package dht

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
MAX_NODE_ENTRY_COUNT  = 10
)

type Prefix struct {
	Size    int
	Flags   []byte
	Content []byte
}

const MaxPaddingSize = 256 // This is an arbitrary value, adjust according to your needs

type IndexEntry struct {
}

type KeySpec map[string]int

type Pht struct {
	KeySpec KeySpec

	Node *Node
}

// Key represents a key in the DHT.
type Key struct {
	Part1 string
    Part2 int
}

// LookupCallback is a callback function type for the lookup result.
type LookupCallback func(vals []IndexEntry, p Prefix)

// LookupCallbackWrapper is a function type for the lookup callback.
type LookupCallbackWrapper func(vals []IndexEntry, p Prefix)

// DoneCallbackSimple is a function type for the done callback.
type DoneCallbackSimple func(ok bool)

// RealInsertCallback is a function type for callbacks used in the insert process.
type RealInsertCallback func(p Prefix, entry IndexEntry)

func NewPrefix() *Prefix {
	return &Prefix{}
}

func NewPrefixWithInfoHash(h InfoHash) *Prefix {
	return &Prefix{
		Size:    len(h) * 8,
		Content: h[:], // Assuming InfoHash is a type that can be directly assigned to Content
	}
}

func NewPrefixWithBlob(d, f []byte) *Prefix {
	return &Prefix{
		Size:    len(d) * 8,
		Flags:   f,
		Content: d,
	}
}

func NewPrefixWithPrefix(p *Prefix, first int) *Prefix {
	size := min(first, len(p.Content)*8)
	newPrefix := &Prefix{
		Size:    size,
		Content: make([]byte, size/8),
	}
	copy(newPrefix.Content, p.Content[:size/8])

	rem := size % 8
	if len(p.Flags) > 0 {
		newPrefix.Flags = make([]byte, len(p.Content[:size/8]))
		copy(newPrefix.Flags, p.Flags[:size/8])
		if rem > 0 {
			newPrefix.Flags = append(newPrefix.Flags, p.Flags[size/8]&(0xFF<<(8-rem)))
		}
	}

	if rem > 0 {
		newPrefix.Content = append(newPrefix.Content, p.Content[size/8]&(0xFF<<(8-rem)))
	}

	return newPrefix
}

func (p *Prefix) GetPrefix(length int) *Prefix {
	if abs(length) >= len(p.Content)*8 {
		panic("length larger than prefix size")
	}
	if length < 0 {
		length += p.Size
	}
	return NewPrefixWithPrefix(p, length)
}

func (p *Prefix) IsFlagActive(pos int) bool {
	if len(p.Flags) == 0 {
		return true
	}
	return isActiveBit(p.Flags, pos)
}

func (p *Prefix) IsContentBitActive(pos int) bool {
	return isActiveBit(p.Content, pos)
}

func (p *Prefix) GetFullSize() *Prefix {
	return NewPrefixWithPrefix(p, len(p.Content)*8)
}
func (p *Prefix) swapContentBit(bit int) {
    if bit >= HashSize*8 {
        panic("bit index out of range")
    }

    byteIndex := bit / 8        // Find the index of the byte containing the bit
    bitPosition := bit % 8      // Find the position of the bit within the byte
    mask := byte(1 << (7 - bitPosition)) // Create a mask to toggle the bit

	h := p.Content
    h[byteIndex] ^= mask // Toggle the bit
}

func (p *Prefix) GetSibling() *Prefix {
	copy := NewPrefixWithPrefix(p, p.Size)
	if p.Size > 0 {
		copy.swapContentBit(p.Size - 1)
	}
	return copy
}

func (p Prefix) String() string {
	return fmt.Sprintf("Prefix:\n\tContent_: \"%s\"\n\tFlags_: \"%s\"\n",
		blobToString(p.Content),
		blobToString(p.Flags))
}

// hash method implementation depends on InfoHash type and its method

func blobToString(blob []byte) string {
	var builder strings.Builder
	for _, b := range blob {
		builder.WriteString(fmt.Sprintf("%08b ", b))
	}
	return builder.String()
}

func CommonPrefixBits(p1, p2 *Prefix) int {
	longestPrefixSize := min(p1.Size, p2.Size)
	for i := 0; i < longestPrefixSize; i++ {
		if p1.Content[i] != p2.Content[i] || !p1.IsFlagActive(i) || !p2.IsFlagActive(i) {
			break
		}
	}

	// Additional implementation...

	return 0 // Placeholder
}

// Additional methods (swapContentBit, swapFlagBit, addPaddingContent, updateFlags)...

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isActiveBit(b []byte, pos int) bool {
	// Implementation...
	return false
}

func (pht *Pht) Insert(kp Prefix, entry IndexEntry, lo, hi *int, timeP time.Time, checkSplit bool, doneCb func(ok bool)) {
	if time.Since(timeP) > 10 {
		return
	}

	// ... Additional logic ...

	// This is a placeholder for the lookupStep function. You need to implement this based on your requirements.
	// pht.lookupStep(kp, lo, hi, func(vals []IndexEntry, prefix Prefix) {
	// 	// RealInsertCallback implementation
	// 	realInsert := func(p Prefix, entry IndexEntry) {
	// 		// Implement the real insert logic here
	// 		// For example, updating a canary, checking the PHT update, inserting into a cache, etc.
	// 	}

	// 	// Logic after lookupStep
	// 	if !checkSplit || len(vals) < MAX_NODE_ENTRY_COUNT {
	// 		// Implement logic based on whether to split or not
	// 	} else {
	// 		// Split logic if needed
	// 	}
	// }, doneCb)
}


func (pht *Pht) lookupStep(p Prefix, lo, hi *int, vals *[]IndexEntry, cb LookupCallbackWrapper, doneCb DoneCallbackSimple, maxCommonPrefixLen *int, start int, allValues bool) {
    // This is a placeholder for your DHT get operation and handling its response.
    // The actual implementation will depend on your DHT's API and logic.

    // This function will likely need to interact with the DHT network,
    // which may involve asynchronous operations. In Go, this is typically
    // handled using goroutines and channels.

    // For example, you might start a goroutine for a DHT get operation:
    go func() {
        // Perform the DHT get operation...
        // On completion, process the results...

        // Call the callback with the results
        cb(*vals, p)

        // Finally, call the done callback
        doneCb(true)
    }()

    // Additional logic for handling lookup steps...
}

func (pht *Pht) Lookup(k Key, cb LookupCallback, doneCb DoneCallbackSimple, exactMatch bool) {
    prefix := pht.Linearize(k) // Assuming Linearize is a method to convert a Key to a Prefix
    values := make([]IndexEntry, 0)

    lo := new(int)
    hi := new(int)
    *lo = 0
    *hi = prefix.Size // Assuming Prefix has a Size field

    var maxCommonPrefixLen *int
    if !exactMatch {
        maxCommonPrefixLen = new(int)
    }

    pht.lookupStep(prefix, lo, hi, &values, func(entries []IndexEntry, p Prefix) {
        // Process the entries as per your logic
        cb(entries, p)
    }, doneCb, maxCommonPrefixLen, -1, true) // Assuming true for allValues for simplicity
}

func (pht *Pht) Linearize(k Key) Prefix {
    // Logic to convert a multi-dimensional key into a unidimensional prefix
    // This is an oversimplified example. You'll need to adjust it to your actual key structure and requirements.

    var allPrefixes []Prefix
    for partName, length := range pht.KeySpec {
        // Assuming 'k' provides a way to get its parts as []byte
        part := k.GetPartAsBytes(partName) // Simplified; replace with actual method to get key part
        paddedPart := padToLength(part, length)
        allPrefixes = append(allPrefixes, Prefix{Content: paddedPart, Size: len(paddedPart)})
    }

    return pht.ZCurve(allPrefixes)
}

func (k Key) GetPartAsBytes(partName string) []byte {
    switch partName {
    case "Part1":
        return []byte(k.Part1)
    case "Part2":
		res, err := intToBytes(k.Part2)
		if err != nil {
			// TODO(nvn): handle the error better
			fmt.Print(err)
		}
        return res
    // handle other parts...
    default:
        return nil // or handle unknown parts as appropriate
    }
}

func padToLength(data []byte, length int) []byte {
    if len(data) >= length {
        return data[:length]
    }
    padding := make([]byte, length-len(data))
    return append(data, padding...)
}

func (pht *Pht) ZCurve(prefixes []Prefix) Prefix {
    // Interleave bits of prefixes. This is a conceptual example.
    // The actual implementation will depend on how you need to interleave the bits.

    var result []byte
    for _, prefix := range prefixes {
        // This is a placeholder logic. You need to implement the actual interleaving algorithm.
        result = append(result, prefix.Content...) // Simplified; replace with actual interleaving logic
    }
    return Prefix{Content: result, Size: len(result)}
}

// intToBytes converts an integer to a byte slice.
// This example uses big-endian format, but you can use binary.LittleEndian if needed.
func intToBytes(n int) ([]byte, error) {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.BigEndian, int64(n))
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func NewPht(ctx context.Context, node *Node) *Pht {
	return &Pht{
		Node: node,
	}
}