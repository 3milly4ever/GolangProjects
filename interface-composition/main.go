package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
)

type HashReader interface {
	io.Reader
	hash() string
}

func main() {
	payload := []byte("hello high value software engineer")
	hashAndBroadcast(NewHashReader(payload))
}

type hashReader struct {
	//these are embedded fields so they don't have explicit names
	*bytes.Reader               //reads in bytes
	buf           *bytes.Buffer //the bytes get stored in this buffer which is like an array of bytes which can expand to as much as it needs to
}

// constructor function for hash reader struct
func NewHashReader(b []byte) *hashReader {
	return &hashReader{
		Reader: bytes.NewReader(b),
		buf:    bytes.NewBuffer(b),
	}
}

func (h *hashReader) hash() string {
	return hex.EncodeToString(h.buf.Bytes())
}

func hashAndBroadcast(r HashReader) error {
	hash := r.hash()
	fmt.Println(hash)
	return broadcast(r)
}

func broadcast(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	fmt.Println("string of the bytes(hashcode): ", string(b))

	return nil
}
