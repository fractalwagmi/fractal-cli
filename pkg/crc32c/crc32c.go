package crc32c

import (
	"hash/crc32"
	"io"
	"os"
)

// Table is a crc32 table based on the Castagnoli polynomial, and can be used
// to compute CRC32-C hashes, which are used on Google Cloud Storage for example.
var table = crc32.MakeTable(crc32.Castagnoli)

func GenerateCrc32C(file string) ([]byte, error) {
	fr, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	hasher := crc32.New(table)
	_, err = io.Copy(hasher, fr)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
