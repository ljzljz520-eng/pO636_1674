package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"repairdesk.local/internal/service"
	"repairdesk.local/internal/store"
	"repairdesk.local/internal/transport"
)

func main() {
	path := flag.String("db", "repairdesk.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	db, e := store.Open(*path)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	svc := service.New(db)
	if e = svc.Seed(); e != nil {
		log.Fatal(e)
	}
	fmt.Println("repairdesk listening on", *addr)
	if e = http.ListenAndServe(*addr, transport.New(svc)); e != nil && !os.IsTimeout(e) {
		log.Fatal(e)
	}
}
