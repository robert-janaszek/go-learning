package course

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ex 15
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("%s %s took %v", r.Method, r.URL.Path, duration)
	})
}

func Day5c() {
	// ex 11
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// ex 12
	http.HandleFunc("GET /api/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		u := user{
			Name: "Mark",
			Addr: address{
				City:    "Warsaw",
				ZipCode: "01-212",
			},
		}
		json.NewEncoder(w).Encode(u)
	})

	// ex 13
	http.HandleFunc("POST /api/user", func(w http.ResponseWriter, r *http.Request) {
		u := user{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			fmt.Println("error ", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			return
		}

		json.NewEncoder(w).Encode(u)
	})

	// ex 14
	http.HandleFunc("GET /api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		w.Write([]byte(id))
	})

	err := http.ListenAndServe(":8080", LoggingMiddleware(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
