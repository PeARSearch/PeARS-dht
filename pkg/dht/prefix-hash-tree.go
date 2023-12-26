package dht

import (
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

type IndexEntry struct {
	// Define the properties of IndexEntry as needed
}

// Define other necessary structs and types...

// Pht represents a DHT.
type Pht struct {
	// Include necessary fields...
}

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