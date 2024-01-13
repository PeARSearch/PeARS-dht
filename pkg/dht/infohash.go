package dht

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"time"
)

const HashSize = 20 // HASH_LEN is 20 bytes

type InfoHash [HashSize]byte

// NewInfoHash creates a new InfoHash from a byte slice.
// It hashes the input if it's longer than HashSize.
func NewInfoHash(data []byte) InfoHash {
	if len(data) == HashSize {
		var h InfoHash
		copy(h[:], data)
		return h
	}
	return sha1.Sum(data)
}

// NewInfoHashFromHex creates a new InfoHash from a hexadecimal string.
func NewInfoHashFromHex(hexStr string) (InfoHash, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return InfoHash{}, err
	}
	return NewInfoHash(bytes), nil
}

// Equals checks if two InfoHashes are equal.
func (h InfoHash) Equals(other InfoHash) bool {
	return h == other
}

// Less checks if one InfoHash is less than another. Useful for sorting.
func (h InfoHash) Less(other InfoHash) bool {
	for i := 0; i < HashSize; i++ {
		if h[i] != other[i] {
			return h[i] < other[i]
		}
	}
	return false
}

// XOR returns the bitwise XOR of two InfoHashes.
func (h InfoHash) XOR(other InfoHash) InfoHash {
	var result InfoHash
	for i := 0; i < HashSize; i++ {
		result[i] = h[i] ^ other[i]
	}
	return result
}

// String returns the hexadecimal string representation of the InfoHash.
func (h InfoHash) String() string {
	return hex.EncodeToString(h[:])
}

// RandomInfoHash generates a random InfoHash.
func RandomInfoHash() InfoHash {
	var h InfoHash
	rand.Seed(time.Now().UnixNano())
	for i := range h {
		h[i] = byte(rand.Intn(256))
	}
	return h
}

// CommonBits returns the number of common leading bits between two InfoHashes.
func CommonBits(h1, h2 InfoHash) int {
	for i := 0; i < HashSize; i++ {
		diff := h1[i] ^ h2[i]
		if diff == 0 {
			continue
		}
		for j := 0; j < 8; j++ {
			if diff&(1<<(7-j)) != 0 {
				return i*8 + j
			}
		}
	}
	return HashSize * 8
}
