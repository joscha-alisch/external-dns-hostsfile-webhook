package main

import (
	"fmt"
	"github.com/joscha-alisch/external-dns-hostsfile-webhook/internal/api"
	"github.com/joscha-alisch/external-dns-hostsfile-webhook/internal/provider"
	"github.com/spf13/pflag"
	"log"
	"net/http"
)

type config struct {
	Filepath     string
	Port         int
	DomainFilter string
}

func main() {
	port := pflag.IntP("port", "p", 8888, "port to listen on")
	path := pflag.StringP("filepath", "f", "./hosts", "path to hosts file")
	pflag.Parse()

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), api.New(provider.New(*path), *path))
	if err != nil {
		log.Fatal(err)
	}
}
