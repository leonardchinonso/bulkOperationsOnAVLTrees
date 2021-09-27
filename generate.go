package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	nodes = flag.Int("nodes", 10000, "number of nodes in the tree before upsert")
	dict  = flag.Int("dict", 100, "number of nodes in the upsert dictionary")
)

// generate create file with given name and number of 4-byte pseudorandom records
func generate(filename string, count int) error {
	var f1 *os.File
	var err error
	if f1, err = os.Create(filename); err != nil {
		return err
	}
	defer f1.Close()
	w1 := bufio.NewWriter(f1)
	defer w1.Flush()
	var buf [4]byte
	for i := 0; i < count; i++ {
		binary.BigEndian.PutUint32(buf[:], rand.Uint32())
		if _, err = w1.Write(buf[:]); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse() // Parse command line flags
	// Seed pseudo random number generator
	rand.Seed(time.Now().UnixNano())
	if err := generate("nodes.dat", *nodes); err != nil {
		fmt.Printf("Error creating nodes.dat: %v\n", err)
	}
	if err := generate("dict.dat", *dict); err != nil {
		fmt.Printf("Error creating dict.dat: %v\n", err)
	}
}
