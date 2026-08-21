package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"independentweeklylog/internal/api"
	"independentweeklylog/internal/archive"
	"independentweeklylog/internal/query"
	"independentweeklylog/internal/repository"
	"independentweeklylog/internal/review"
	"independentweeklylog/internal/store"
	"independentweeklylog/internal/workflow20"
)

func main() {
	path := flag.String("db", "weeklylog.db", "bbolt database path")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	db, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := repository.New(db)
	workflow := workflow20.New(repo)
	reviewer := review.New(repo)
	finder := query.New(repo)
	archiver := archive.New(repo)
	server := api.NewServer(workflow, reviewer, finder)
	handler := server.Handler()
	if mux, ok := handler.(*http.ServeMux); ok {
		mux.Handle("/archive", api.ArchiveHandler(repo, archiver))
	}
	fmt.Printf("independent game weekly log listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, handler))
}
