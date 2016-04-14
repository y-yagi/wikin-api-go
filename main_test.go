package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/zenazn/goji/web"
)

func TestMain(m *testing.M) {
	db, _ = gorm.Open("postgres", "user=yaginuma dbname=wikin_test")
	code := m.Run()
	defer os.Exit(code)
	db.Delete(Page{})
}

func ParseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}

func Test_Page(t *testing.T) {

	m := web.New()
	Route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	page := Page{
		Id:    1,
		Title: "test title",
		Body:  "test body",
	}

	result := db.Create(&page)
	if result.Error != nil {
		t.Error("test data create error", result.Error)
	}

	res, err := http.Get(ts.URL + "/pages/1")
	if err != nil {
		t.Error("unexpected")
	}
	c, s := ParseResponse(res)
	if s != http.StatusOK {
		t.Error("invalid status code", s)
	}

	dec := json.NewDecoder(strings.NewReader(c))
	var responsePage Page
	dec.Decode(&responsePage)

	if page.Id != 1 {
		t.Error("invalid id: ", page.Id)
	}
	if page.Title != "test title" {
		t.Error("invalid title: ", page.Title)
	}
	if page.Body != "test body" {
		t.Error("invalid body: ", page.Body)
	}

}
