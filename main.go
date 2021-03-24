package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var concurrency int

func init() {
	flag.IntVar(&concurrency, "concurrency", 1, "how many to run at once")
	flag.Parse()
}

func main() {
	file, err := os.Open("./manifest.yml")
	if err != nil {
		log.Fatalf("error opening manifest: %v", err)
	}
	defer file.Close()

	yamlManifest, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading manifest: %v", err)
	}

	manifest := &Manifest{}

	yaml.Unmarshal(yamlManifest, manifest)

	manifest.Run(concurrency)

}
