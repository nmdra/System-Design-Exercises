package main

import (
	"log"
	"net/http"
	"time"
)

func process() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Background Service Start.")
			time.Sleep(10 * time.Second)
			log.Println("Background Service Stop.")
		}()
		w.Write([]byte("Done\n"))
	})
}
