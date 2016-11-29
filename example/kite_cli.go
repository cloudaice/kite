package main

import (
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

var (
	costs []float64
	lock  sync.Mutex
)

func main() {
	var length, num int
	flag.IntVar(&length, "l", 256, "msg length for testing")
	flag.IntVar(&num, "n", 1, "try count for testing")
	flag.Parse()
	wg := sync.WaitGroup{}
	t0 := time.Now()

	for i := 0; i < num; i++ {
		wg.Add(1)

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
	testMsg := kite.RandBytes(length)
	conn, err := net.Dial("tcp", "127.0.0.1:9876")
	if err != nil {
		fmt.Printf("Client Dial Error: %s\n", err)
		return
	}
	defer conn.Close()

	cli := kite.NewClient(conn)
	err = cli.SendHandshake()
	if err != nil {
		fmt.Printf("Client SendHandshake Error: %s\n", err)
		return
	}

	for i := 0; i < 100000; i++ {
		err = cli.SendData(testMsg, kite.D_NONE)
		if err != nil {
			fmt.Printf("Client SendData Error: %s\n", err)
			return
		}
		_, err = cli.ReadPackage()
		if err != nil {
			fmt.Printf("ReadPackage Error: %s\n", err)
			return
		}
	}

	err = cli.SendClose()
	if err != nil {
		fmt.Printf("Client SendClose Error: %s\n", err)
		return
	}
	err = cli.ReadCloseAck()
	if err != nil {
		fmt.Printf("Client ReadCloseAck Error: %s\n", err)
		return
	}
}
