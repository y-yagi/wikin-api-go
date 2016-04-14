package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/zenazn/goji/web"
)

func ParseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}

func Test_Page(t *testing.T) {
	db, _ = gorm.Open("postgres", "user=yaginuma dbname=wikin_test")

	m := web.New()
	Route(m)
	ts := httptest.NewServer(m)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/pages/2")
	if err != nil {
		t.Error("unexpected")
	}
	c, s := ParseResponse(res)
	if s != http.StatusOK {
		t.Error("invalid status code")
	}
	if c != `{"id":2,"title":"c","body":"testmessage","created_at":"2016-04-07T05:38:22.798216Z","updated_at":"2016-04-09T17:40:04.128029Z"}` {
		t.Error("invalid response")
	}
}
