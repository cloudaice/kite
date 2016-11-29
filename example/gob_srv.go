package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type P string

func (p0 *P) equal(p1 *P) bool {
	return string(*p0) == string(*p1)
}

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:9999", nil)
	}()

	l, err := net.Listen("tcp", ":19876")
	if err != nil {
		fmt.Printf("Server Listen Error: %s\n", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Server Connect Error: %s\n", err)
		}
		go handleConnetion(conn)
	}
}

func handleConnetion(conn net.Conn) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	var err error

	var p P
	pc := P("Close")

	t0 := time.Now()
	fmt.Printf("[%s] Client %s Conneted.\n", t0, conn.RemoteAddr().String())

	for ok := false; !ok; {
		err = dec.Decode(&p)
		if err != nil {
			fmt.Printf("Decode Error: %s\n", err)
			break
		}
		if p.equal(&pc) {
			ok = true
		}
		err = enc.Encode(p)
		if err != nil {
			fmt.Printf("Encode Error: %s\n", err)
		}
	}
	c := time.Since(t0).Seconds()
	fmt.Printf("[%s] Client %s Close! - Costs %f\n", time.Now(), conn.RemoteAddr().String(), c)
}
