package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"kite"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type P string

func (p0 *P) equal(p1 *P) bool {
	return string(*p0) == string(*p1)
}

func main() {
	var length, num int
	flag.IntVar(&length, "l", 256, "msg length for testing")
	flag.IntVar(&num, "n", 1, "try count for testing")
	flag.Parse()
	wg := sync.WaitGroup{}
	t0 := time.Now()
	for i := 0; i < num; i++ {
		wg.Add(1)
		//fmt.Printf("Connection %d\n", i)
		go func() {
			test(length)
			wg.Done()
		}()
	}
	wg.Wait()
	duration := time.Since(t0).Seconds()
	fmt.Printf("%d TPS\n", int(float64(100000*num)/duration))
}

func test(length int) {
	var p0 P
	p := P(string(kite.RandBytes(length)))
	pc := P("Close")

	conn, err := net.Dial("tcp", "127.0.0.1:19876")
	if err != nil {
		fmt.Printf("Client Dial Error: %s\n", err)
	}
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	for i := 0; i < 100000; i++ {
		err = enc.Encode(p)
		if err != nil {
			fmt.Printf("Encode Error: %s\n", err)
			break
		}
		err = dec.Decode(&p0)
		if err != nil {
			fmt.Printf("Decode Error: %s\n", err)
			break
		}
	}

	err = enc.Encode(&pc)
	if err != nil {
		fmt.Printf("Encode Error: %s\n", err)
	}

	err = dec.Decode(&p0)
	if err != nil {
		fmt.Printf("Decode Error: %s\n", err)
	}
	if !p0.equal(&pc) {
		fmt.Println("Not Close....")
	}
}
