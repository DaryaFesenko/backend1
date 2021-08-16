package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go watchSignal(cancel)

	d := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}
	conn, err := d.DialContext(ctx, "tcp", "127.127.127.127:9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(io.Copy(os.Stdout, conn))
}

func watchSignal(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, syscall.SIGINT)

	<-osSignalChan

	log.Println("user interrupted")
	cancel()
}
