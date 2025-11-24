package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"sync"
	"time"

	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

var n = flag.Int("n", 100000, "number of trials")
var t = flag.Int("t", 1, "number of threads")
var b = flag.Int("b", 100000, "batch per task")
var r = flag.Int("r", 1000000, "report interval")
var v = flag.Int("v", 3, "verbosity")
var hv = flag.Int("h", 3, "hash version")

func main() {
	flag.Parse()

	start := time.Now()

	prefix := []byte("data")

	var wg sync.WaitGroup
	var m sync.Mutex

	tid := 0

	if *hv == 3 {
		fmt.Printf("use Keccak256\n")
		fmt.Println("keccah256 test vector")
		k := sha3.NewLegacyKeccak256()
		k.Write([]byte{})
		h := k.Sum(nil)
		fmt.Println(hex.EncodeToString(h))
	} else if *hv == 2 {
		fmt.Printf("use Sha256\n")
	} else {
		fmt.Println("unsupported hash")
		return
	}

	for i := 0; i < *t; i++ {
		wg.Add(1)

		go func(thread int) {
			defer wg.Done()

			for {
				m.Lock()
				ltid := tid
				tid = tid + *b
				m.Unlock()

				if ltid >= *n {
					break
				}

				if *v > 3 {
					fmt.Printf("thread %d: %d to %d\n", thread, ltid, ltid+*b)
				}

				if *v >= 3 && ltid%*r == 0 {
					elapsed := time.Since(start)
					fmt.Printf("used time %f, hps %f\n", elapsed.Seconds(), float64(ltid)/elapsed.Seconds())
				}

				for j := ltid; j < ltid+*b; j++ {
					if *hv == 3 || *hv == 2 {
						var k hash.Hash
						if *hv == 3 {
							k = sha3.NewLegacyKeccak256()
						} else {
							k = sha256.New()
						}
						buf := make([]byte, 4096+32)
						binary.BigEndian.PutUint64(buf[4096:], uint64(j))
						k.Write(prefix)
						k.Write(buf)
						k.Sum(nil)
					} else {
						panic("unsupported hash methods")
					}
				}
			}

		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("used time %f, hps %f\n", elapsed.Seconds(), float64(*n)/elapsed.Seconds())
}
