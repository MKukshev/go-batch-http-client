package main

import (
	"fmt"
	"go-batch-http-client/app"
	. "go-batch-http-client/model"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	f, err := os.Open("config/config.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable configuration file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration file currupted: %v\n", err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(cfg.Logger.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	app.Run(cfg)
}
