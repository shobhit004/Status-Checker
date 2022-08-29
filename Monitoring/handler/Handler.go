package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var URLs = make(map[string]string)

func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		param := r.URL.Query().Get("name")
		if param != "" {
			status, _ := URLs[param]
			fmt.Fprintln(w, fmt.Sprintf("Status of %s: %s", param, status))

		} else {
			for key, val := range URLs {
				fmt.Fprintln(w, key, " ", val)
			}
		}
	case "POST":
		var websites []string
		err := json.NewDecoder(r.Body).Decode(&websites)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Error:%+v", err))
		}

		c := make(chan string)
		for _, url := range websites {
			go checkStatus(url, c)
		}
		for link := range c {
			go func(l string) {
				time.Sleep(5 * time.Second)
				checkStatus(l, c)
			}(link)
		}
	default:
		fmt.Fprint(w, "Sorry, only GET and POST methods are supported.")
	}
}
func checkStatus(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, " is not working")
		URLs[link] = "not working"
		c <- link
		return
	}
	fmt.Println(link, " works")
	URLs[link] = "working"
	c <- link
	return
}
