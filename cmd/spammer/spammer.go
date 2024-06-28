package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	network = flag.String("network", "localhost", "Network")
	minPort = flag.Int("minPort", 0, "port start range value")
	maxPort = flag.Int("maxPort", 65535, "port end range value")
	payload = flag.String("payload", "spammer", "Spammer payload")

	random   = flag.Bool("random", false, "Random payload")
	interval = flag.Int("interval", 1000, "Interval in millis")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func spammer(network string, port int) error {
	hostport := net.JoinHostPort(network, fmt.Sprintf("%d", port))
	conn, err := net.Dial("udp4", hostport)
	if err != nil {
		return err
	}

	fmt.Println("initialized spammer:", hostport)

	b := make([]byte, 255)
	data := func() []byte {
		if *random {
			rand.Read(b)
			return b
		} else {
			return []byte(*payload)
		}

	}
	for {
		n, err := conn.Write(data())
		if err != nil {
			fmt.Println(hostport, "errored:", err)
			continue
		}

		fmt.Println(hostport, "written", n)

		time.Sleep(time.Duration(*interval) * time.Millisecond)
	}
}

func run() error {

	var wg sync.WaitGroup

	for port := *minPort; port <= *maxPort; port++ {
		wg.Add(1)
		port := port
		go func() {
			defer wg.Done()
			if err := spammer(*network, port); err != nil {
				fmt.Println("spammer", port, "is down:", err)
			}
		}()
	}

	wg.Wait()

	return nil
}
