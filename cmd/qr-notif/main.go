package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	h "github.com/alfianX/jserver/helper"
	qrnotif "github.com/alfianX/jserver/internal/qr-notif"
)

var (
	version     = "1.0.0" // Application version
	showVersion = flag.Bool("version", false, "Display the application version")
)

func main() {
	flag.Parse()

	// If "--version" flag is provided, display version and exit
	if *showVersion {
		fmt.Printf("App Version: %s\n", version)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		run(context.Background())
	}()

	wg.Wait()
}

func run(ctx context.Context) {
	server, err := qrnotif.NewService()
	if err != nil {
		h.ErrorLog(fmt.Sprintf("%+v", err))
		log.Fatalf("%+v", err)
	}

	err = server.Run(ctx)
	if err != nil {
		h.ErrorLog(fmt.Sprintf("%+v", err))
		log.Fatalf("%+v", err)
	}
}
