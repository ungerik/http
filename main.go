package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:\nhttp get|head|post <URL> [<POST_ARG>=<VALUE>]*")
		os.Exit(0)
	}

	method := strings.ToUpper(os.Args[1])
	URL := os.Args[2]
	if !strings.HasPrefix(URL, "http") {
		URL = "http://" + URL
	}

	switch method {
	case "HEAD":
		r, err := http.Head(URL)
		if err != nil {
			panic(err)
		}
		r.Header.Write(os.Stdout)

	case "GET":
		r, err := http.Get(URL)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		io.Copy(os.Stdout, r.Body)

	case "POST":
		var values url.Values
		for _, arg := range os.Args[3:] {
			kv := strings.Split(arg, "=")
			if len(kv) != 2 {
				panic("Post values must be of form key=value, got " + arg)
			}
			key := strings.Trim(strings.TrimSpace(kv[0]), `"'`)
			value := strings.Trim(strings.TrimSpace(kv[1]), `"'`)
			values.Add(key, value)
		}
		r, err := http.PostForm(URL, values)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		io.Copy(os.Stdout, r.Body)

	default:
		panic("Invalid HTTP method " + method)
	}
}
