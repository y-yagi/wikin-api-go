package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

var (
	db *gorm.DB
)

func getPages(c echo.Context) error {
	var pageSize int

	if c.QueryParam("recent_pages") != "" {
		pageSize = 10
	} else {
		pageSize = 100
	}

	var pages []Page = make([]Page, pageSize)
	var buffer bytes.Buffer
	db.Order("updated_at DESC").Limit(pageSize).Find(&pages)

	mappage, _ := json.Marshal(pages)
	buffer.WriteString(string(mappage))

	return c.String(http.StatusOK, buffer.String())
}

func getPage(c echo.Context) error {
	var page Page

	db.Find(&page, c.Param("id"))
	if page.Id == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	mappage, _ := json.Marshal(page)
	return c.String(http.StatusOK, string(mappage))
}

func searchPages(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return c.String(http.StatusOK, "")
	}

	var buffer bytes.Buffer
	pageSearcher := PageSearcher{query, db}
	pages, err := pageSearcher.Matches()

	if err != nil {
		// TODO: error handling
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if len(pages) == 0 {
		return c.String(http.StatusOK, "")
	}

	mappage, _ := json.Marshal(pages)
	buffer.WriteString(string(mappage))

	return c.String(http.StatusOK, buffer.String())
}

func updatePage(c echo.Context) error {
	var page Page

	db.Find(&page, c.Param("id"), ".json", "", -1)
	if page.Id == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	page.Body = c.FormValue("page[body]")
	db.Save(&page)
	mappage, _ := json.Marshal("status=ok")
	return c.String(http.StatusOK, string(mappage))
}

func Route(e *echo.Echo) {
	e.GET("/pages", getPages)
	e.GET("/pages/search", searchPages)
	e.GET("/pages/:id", getPage)
	e.PATCH("/pages:id", updatePage)
}

func main() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	port := os.Getenv("PORT")
	e := echo.New()

	if os.Getenv("BASIC_AUTH_USER") != "" && os.Getenv("BASIC_AUTH_PASSWORD") != "" {
		e.Use(middleware.BasicAuth(func(username, password string) bool {
			if username == os.Getenv("BASIC_AUTH_USER") && password == os.Getenv("BASIC_AUTH_PASSWORD") {
				return true
			}
			return false
		}))
	}
	Route(e)
	e.Logger.Fatal(e.Start(":" + port))
}
