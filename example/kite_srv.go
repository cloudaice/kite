package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"kite"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:9999", nil)
	}()

	l, err := net.Listen("tcp", ":9876")
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
	srv := kite.NewServer(conn)
	err := srv.WaitHandshake()
	if err != nil {
		fmt.Printf("WaitHandshake Error: %s\n", err)
		return
	}
	t0 := time.Now()
	fmt.Printf("[%s] Client %s Conneted.\n", t0, conn.RemoteAddr().String())

	for ok := false; !ok; {
		pp, err := srv.ReadPackage()
		if err != nil {
			fmt.Printf("Server ReadPackage Error: %s\n", err)
			break
		}
		switch pp.PackageType() {
		case kite.T_CLOSE:
			err = srv.SendCloseAck()
			if err != nil {
				fmt.Printf("Server SendCloseAck Error: %s\n", err)
				break
			}
			c := time.Since(t0).Seconds()
			fmt.Printf("[%s] Client %s Close! - Costs %f\n", time.Now(), conn.RemoteAddr().String(), c)
			ok = true
			break
		default:
			err = srv.SendData(pp.Content(), kite.D_NONE)
			if err != nil {
				fmt.Printf("Server SendData Error: %s\n", err)
			}
		}
	}
}
