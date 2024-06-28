package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	network = flag.String("network", "", "Network to sniff traffic")
	minPort = flag.Int("minPort", 0, "port start range value")
	maxPort = flag.Int("maxPort", 65535, "port end range value")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func sniffer(network string, port string) error {
	hostport := net.JoinHostPort(network, port)

	pc, err := net.ListenPacket("udp4", hostport)
	if err != nil {
		return err
	}

	fmt.Println("initialized sniffer", hostport)

	buf := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println(hostport, "errored: ", err)
			continue
		}

		fmt.Println(hostport, "read packet", n, "from:", addr.String())
		fmt.Println("--- ASCII ---")
		fmt.Println(string(buf[:n]))
		fmt.Println("--- HEX ---")
		fmt.Println(hex.EncodeToString(buf[:n]))

		time.Sleep(100 * time.Millisecond)
	}
}

func run() error {

	var wg sync.WaitGroup

	for port := *minPort; port <= *maxPort; port++ {
		wg.Add(1)
		port := port
		go func() {
			defer wg.Done()
			sniffer(*network, fmt.Sprintf("%d", port))
		}()
	}

	wg.Wait()

	return nil
}
