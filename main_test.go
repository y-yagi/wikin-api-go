package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/zenazn/goji/web"
)

func startServer() *httptest.Server {
	m := web.New()
	Route(m)
	ts := httptest.NewServer(m)
	return ts
}

func parseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}

func generateTestData() {
	db.Delete(Page{})
	var page Page
	for i := 1; i < 6; i++ {
		page = Page{
			Id:    i,
			Title: "test title" + strconv.Itoa(i),
			Body:  "test body" + strconv.Itoa(i),
		}

		result := db.Create(&page)
		if result.Error != nil {
			fmt.Printf("test data create error %s", result.Error)
		}
	}
}

func TestMain(m *testing.M) {
	db, _ = gorm.Open("postgres", "dbname=wikin_test  sslmode=disable")
	defer db.Close()
	generateTestData()
	code := m.Run()
	defer os.Exit(code)
}

func Test_Page(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	res, err := http.Get(ts.URL + "/pages/1")
	if err != nil {
		t.Error("unexpected")
	}
	c, s := parseResponse(res)
	if s != http.StatusOK {
		t.Error("invalid status code", s)
	}

	dec := json.NewDecoder(strings.NewReader(c))
	var page Page
	dec.Decode(&page)

	if page.Id != 1 {
		t.Error("invalid id: ", page.Id)
	}
	if page.Title != "test title1" {
		t.Error("invalid title: ", page.Title)
	}
	if page.Body != "test body1" {
		t.Error("invalid body: ", page.Body)
	}

}

func Test_Pages(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	res, err := http.Get(ts.URL + "/pages")
	if err != nil {
		t.Error("unexpected")
	}
	c, s := parseResponse(res)
	if s != http.StatusOK {
		t.Error("invalid status code", s)
	}

	dec := json.NewDecoder(strings.NewReader(c))
	var pages []Page
	dec.Decode(&pages)

	if len(pages) != 5 {
		t.Error("invalid response: ", c)
	}

	for i, j := 0, 5; i < 5; i++ {
		if pages[i].Id != j {
			t.Error("invalid id: ", pages[i].Id)
		}
		if pages[i].Title != "test title"+strconv.Itoa(j) {
			t.Error("invalid title: ", pages[i].Title)
		}
		if pages[i].Body != "test body"+strconv.Itoa(j) {
			t.Error("invalid body: ", pages[i].Body)
		}
		j--
	}
}

func Test_SearchPages(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	res, err := http.Get(ts.URL + "/pages/search?query=title1")
	if err != nil {
		t.Error("unexpected")
	}
	c, s := parseResponse(res)
	if s != http.StatusOK {
		t.Error("invalid status code", s)
	}

	dec := json.NewDecoder(strings.NewReader(c))
	var pages []Page
	dec.Decode(&pages)

	if len(pages) != 1 {
		t.Error("invalid response: ", c)
	}

	if pages[0].Id != 1 {
		t.Error("invalid id: ", pages[0].Id)
	}
	if pages[0].Title != "test title1" {
		t.Error("invalid title: ", pages[0].Title)
	}
	if pages[0].Body != "test body1" {
		t.Error("invalid body: ", pages[0].Body)
	}

}
