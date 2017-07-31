package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	c := make(chan int64)
	s := make(chan int64)
	l := make(chan int)

	go func() {
		for {
			time.Sleep(time.Second)
			s <- time.Now().Add(time.Minute * -1).Unix()
		}
	}()

	go func() {
		var timestamps []int64
		for {
			select {
			case t := <-c:
				timestamps = append(timestamps, t)
				l <- len(timestamps)
			case t := <-s:
				for i := range timestamps {
					if timestamps[i] > t {
						timestamps = timestamps[i:]
						break
					}
				}
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c <- time.Now().Unix()
		fmt.Fprint(w, <-l)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, nil)
}
