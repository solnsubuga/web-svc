package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/solnsubuga/web-svc/service"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	ttl := flag.Duration("ttl", time.Second*15, "Service TTL check duration")

	svc, err := service.New(*ttl)
	if err != nil {
		log.Fatal(err)
	}

	if err = svc.RegisterSvc(); err != nil {
		log.Fatal(err)
	}

	go svc.UpdateConsul(svc.Check)

	http.Handle("/", svc)

	l := fmt.Sprintf(":%d", *port)
	log.Print("Listening on ", l)
	log.Fatal(http.ListenAndServe(l, nil))
}
