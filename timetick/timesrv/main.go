package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go watchSignal(cancel)

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}
	l, err := cfg.Listen(ctx, "tcp", "127.127.127.127:9000")
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	log.Println("im started!")

	go func() {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		} else {
			wg.Add(1)
			go handleConn(ctx, conn, wg)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("done")
			l.Close()
			wg.Wait()
			log.Println("exit")
			return
		}
	}
}

func watchSignal(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, syscall.SIGINT)

	<-osSignalChan

	log.Println("user interrupted")
	cancel()
}

func handleConn(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	go writeMessageFromConsole(conn)

	// каждую 1 секунду отправлять клиентам текущее время сервера
	tck := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tck.C:
			fmt.Fprintf(conn, "now: %s\n", t)
		}
	}
}

func writeMessageFromConsole(conn net.Conn) {
	_, err := io.Copy(conn, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
