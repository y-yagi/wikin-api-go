package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

type PageSearcher struct {
	query string
	db    *gorm.DB
}

func (searcher PageSearcher) Matches() ([]Page, error) {
	var parts []string
	var bindings []string

	for _, word := range splitQuery(searcher.query) {
		escapedWord := escapeForLikeQuery(word)
		parts = append(parts, "title LIKE $1 OR body LIKE $2")
		bindings = append(bindings, escapedWord)
		bindings = append(bindings, escapedWord)
	}

	var pages []Page
	var page Page

	// TODO: can use multi parameters
	rows, err := db.DB().Query("SELECT * FROM pages WHERE "+strings.Join(parts, " OR "), bindings[0], bindings[1])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&page.Id, &page.Title, &page.Body, &page.Parent_Id, &page.Created_at, &page.Updated_at); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println(pages)
	return pages, nil
}

func splitQuery(query string) []string {
	return regexp.MustCompile("%s+").Split(query, -1)
}

func escapeForLikeQuery(phrase string) string {
	return "%" + phrase + "%"
}
