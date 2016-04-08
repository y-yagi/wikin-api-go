package main

import (
	"fmt"
	"net/http"
	"time"
)

func RecordLog(r *http.Request) {
	fmt.Printf("Started %s %s for %s at %s\n", r.Method, r.RequestURI, r.RemoteAddr, time.Now().Format(time.RFC3339))
}
