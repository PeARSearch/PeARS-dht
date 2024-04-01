package dht

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// TODO(nvn): Make this part of the configuration
	DefaultOrchardIndexerCallerURL = "http://localhost:8080/indexer/hash"
	DefaultOrchardSearchCallerURL = "http://localhost:8080/hash"
	ErrNoSuccessor          = errors.New("cannot find successor")
	ErrNodeExists           = errors.New("node with id already exists")
	ErrKeyNotFound          = errors.New("key not found")
	testURLs                = []string{
		"https://en.wikipedia.org/wiki/Dog",
		"https://en.wikipedia.org/wiki/Cat",
		"https://en.wikipedia.org/wiki/Computer",
		"https://en.wikipedia.org/wiki/Democracy",
		"https://en.wikipedia.org/wiki/Volkswagen",
		"https://en.wikipedia.org/wiki/Cow",
		"https://en.wikipedia.org/wiki/Train",
	}
)

func isEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func isPowerOfTwo(num int) bool {
	return (num != 0) && ((num & (num - 1)) == 0)
}

func randStabilize(min, max time.Duration) time.Duration {
	r := rand.Float64()
	return time.Duration((r * float64(max-min)) + float64(min))
}

// check if key is between a and b, right inclusive
func betweenRightIncl(key, a, b []byte) bool {
	return between(key, a, b) || bytes.Equal(key, b)
}

// Checks if a key is STRICTLY between two ID's exclusively
func between(key, a, b []byte) bool {
	switch bytes.Compare(a, b) {
	case 1:
		return bytes.Compare(a, key) == -1 || bytes.Compare(b, key) >= 0
	case -1:
		return bytes.Compare(a, key) == -1 && bytes.Compare(b, key) >= 0
	case 0:
		return !bytes.Equal(a, key)
	}
	return false
}

// For testing
func GetHashID(key string) []byte {
	h := sha1.New()
	if _, err := h.Write([]byte(key)); err != nil {
		return nil
	}
	val := h.Sum(nil)
	return val
}

func PopulateTestData(ctx context.Context, node *Node) error {
	for _, docUrl := range testURLs {
		infoHash, err := GenerateInfoHash(ctx, DefaultOrchardIndexerCallerURL, map[string]string{"url": docUrl})
		if err != nil {
			return err
		}

		err = node.Set(infoHash, docUrl)
		if err != nil {
			return fmt.Errorf("failed storing infohash %s for url %s: %v", infoHash, docUrl, err)
		}
	}

	return nil
}

func GenerateInfoHash(ctx context.Context, orchardURL string, payload map[string]string) (string, error) {
	// Parse the target URL
	parsedURL, err := url.Parse(orchardURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	// Add 'url' as a query parameter
	queryParams := parsedURL.Query()
	for key, value := range payload {
		queryParams.Set(key, value)
	}
	parsedURL.RawQuery = queryParams.Encode()

	// Make the GET request
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return "", fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to pull hash from orchard, status code: %d, error: %s ", resp.StatusCode, resp.Body)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	infoHash := hashString(string(body))
	logrus.WithField("Info Hash:", infoHash).WithField("payload", payload).Info("got infohash for payload")

	return infoHash, nil
}

func hashString(input string) string {
	// Convert the input string into a byte slice since the hashing function expects bytes
	data := []byte(input)
	// Sum256 returns the SHA256 checksum of the data
	hash := sha256.Sum256(data)
	// Encode the checksum to a hex string
	hexString := hex.EncodeToString(hash[:])
	return hexString
}
