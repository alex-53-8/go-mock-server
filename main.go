package main

import (
	"context"
	"flag"
	"log"
	"mock-server/server"
	"os"
	"os/signal"
	"time"
)

func main() {
	file := flag.String("file", "server.yml", "file with a server model")
	printModel := flag.Bool("print-model", false, "print model")
	cachingEnabled := flag.Bool("caching-enabled", true, "enable caching for responses read from a storage")
	cacheItemMaxSize := flag.Int64("cached-item-max-size", 1024*1024, "max size of one item can be stored in cache")

	flag.Parse()

	cfg := server.Cfg{CacheItemMaxSize: *cacheItemMaxSize, CachingEnabled: *cachingEnabled}

	if file == nil {
		log.Println("filename is nil")
		return
	}
	log.Println("reading a model from a file", *file)

	model, err := server.ReadModel(*file)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if *printModel {
		log.Printf("%+v\n", model)
	}

	srv := server.NewServer(model, cfg)
	go func() {
		err = srv.Listen()
		log.Println(err)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
