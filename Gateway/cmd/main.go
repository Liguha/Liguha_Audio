package main

import (
	"Gateway/internal/app"
	"Gateway/internal/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func main() {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Open(filepath.Join(mydir, "cmd", "/config.yaml"))
	if err != nil {
		log.Fatalf("config read err %v", err)
	}
	defer f.Close()

	var cfg config.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("config parse err %v", err)
	}

	application := app.New(&cfg)
	application.Run()
}
