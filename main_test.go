package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

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
	e := echo.New()

	req, err := http.NewRequest(echo.GET, "/pages/1", nil)
	if err != nil {
		t.Error("unexpected")
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/pages/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	getPage(c)

	if rec.Code != http.StatusOK {
		t.Error("invalid status code", rec.Code)
	}

	dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
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
	e := echo.New()

	req, err := http.NewRequest(echo.GET, "/pages", nil)
	if err != nil {
		t.Error("unexpected")
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	getPages(c)

	if rec.Code != http.StatusOK {
		t.Error("invalid status code", rec.Code)
	}

	dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
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
	e := echo.New()

	req, err := http.NewRequest(echo.GET, "/pages/search?query=title1", nil)
	if err != nil {
		t.Error("unexpected")
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/pages/search")
	c.SetParamNames("query")
	c.SetParamValues("title1")

	searchPages(c)

	dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
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

func Test_UpdatePage(t *testing.T) {
	e := echo.New()

	var bodyStr = []byte("page[body]=update body")
	req, err := http.NewRequest(echo.PATCH, "/pages/1", bytes.NewBuffer(bodyStr))
	if err != nil {
		t.Error("unexpected", err)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/pages/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	updatePage(c)

	if rec.Code != http.StatusOK {
		t.Error("invalid status code", rec.Code)
	}

	var page Page
	db.Find(&page, "1")

	if page.Body != "update body" {
		t.Error("invalid body: ", page.Body)
	}
}
