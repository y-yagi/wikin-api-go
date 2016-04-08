package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var (
	db *gorm.DB
)

func getPages(db *gorm.DB) ([]Page, error) {
	var pages []Page = make([]Page, 100)

	db.Find(&pages)

	return pages, nil
}

func Getpages(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Started %s %s for %s at %s\n", r.Method, r.RequestURI, r.RemoteAddr, time.Now().Format(time.RFC3339))

	var buffer bytes.Buffer

	pages, err := getPages(db)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, page := range pages {
		mappage, _ := json.Marshal(page)
		buffer.WriteString(string(mappage))
	}

	fmt.Fprint(w, buffer.String())
}

func Getpage(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Started %s %s for %s at %s\n", r.Method, r.RequestURI, r.RemoteAddr, time.Now().Format(time.RFC3339))

	var page Page

	// really safe?
	db.Where("id= ?", c.URLParams["id"]).First(&page)
	mappage, _ := json.Marshal(page)
	fmt.Fprint(w, string(mappage))
}

func main() {
	var err error
	db, err = gorm.Open("postgres", "user=yaginuma dbname=wikin")
	port := os.Getenv("PORT")

	if err != nil {
		log.Fatal(err)
		return
	}

	goji.Get("/pages", Getpages)
	goji.Get("/pages/:id", Getpage)
	flag.Set("bind", ":"+port)
	goji.Serve()
}
