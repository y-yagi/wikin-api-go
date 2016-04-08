package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var (
	db *gorm.DB
)

func getPagesFromDB(db *gorm.DB, size int) ([]Page, error) {
	var pages []Page = make([]Page, size)

	db.Order("updated_at DESC").Limit(size).Find(&pages)

	return pages, nil
}

func GetPages(w http.ResponseWriter, r *http.Request) {
	var pageSize int

	if r.URL.Query().Get("recent_pages") != "" {
		pageSize = 10
	} else {
		pageSize = 100
	}

	var buffer bytes.Buffer

	pages, err := getPagesFromDB(db, pageSize)
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

func GetPage(c web.C, w http.ResponseWriter, r *http.Request) {
	var page Page

	// really safe?
	db.Where("id= ?", c.URLParams["id"]).First(&page)
	mappage, _ := json.Marshal(page)
	fmt.Fprint(w, string(mappage))
}

func main() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	port := os.Getenv("PORT")

	if err != nil {
		log.Fatal(err)
		return
	}

	goji.Get("/pages", GetPages)
	goji.Get("/pages/:id", GetPage)
	flag.Set("bind", ":"+port)
	goji.Serve()
}
