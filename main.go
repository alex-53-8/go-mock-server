package main

import (
	"flag"
	"log"
	"mock-server/server"
)

func main() {
	file := flag.String("file", "model.yml", "file with a server model")
	printModel := flag.Bool("print-model", false, "print model")
	flag.Parse()

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

	srv := server.NewServer(model)
	err = srv.Listen()
	log.Println(err)
}
