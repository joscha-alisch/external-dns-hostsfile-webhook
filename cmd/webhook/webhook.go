package main

import (
	"github.com/joscha-alisch/external-dns-hostsfile-webhook/internal/api"
	"github.com/joscha-alisch/external-dns-hostsfile-webhook/internal/provider"
	"log"
	"net/http"
)

func main() {
	path := "./hosts"
	err := http.ListenAndServe(":8888", api.New(provider.New(path), path))
	if err != nil {
		log.Fatal(err)
	}
}