package main

import (
	"backend1/api/handler"
	"backend1/api/server"
	"backend1/app/services/upload"
	"backend1/app/starter"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watchSignal(cancel)

	a := starter.NewApp()
	us := upload.NewUploadService("/home/d/projects/gb/backend1/upload")
	h := handler.NewRouter(us)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}

func watchSignal(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt)

	<-osSignalChan

	log.Println("user interrupted")
	cancel()
}
