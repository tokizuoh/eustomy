package main

import (
	"cdb"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getSameVowelsWords(s string) ([]cdb.RomanWord, error) {
	rws, err := cdb.GetSameVowelsWords("./debug_roman_vowels_20210613_0102.csv", s)
	if err != nil {
		return nil, err
	}

	return rws, nil
}

func checkAuth(r *http.Request) bool {
	reqClientID, reqClientSecret, ok := r.BasicAuth()
	if ok != true {
		return false
	}

	if err := godotenv.Load(); err != nil {
		return false
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	return clientID == reqClientID && clientSecret == reqClientSecret
}

func main() {
	d := flag.Bool("d", false, "debug flag")
	flag.Parse()
	if *d {
		if err := cdb.GenerateCustomDB(); err != nil {
			log.Fatal(err)
		}
		return
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		log.Println(1)
		t := r.URL.Query().Get("target")
		if t == "" {
			// クエリパラメータが不正な値のとき
			log.Fatal(1)
			return
		}

		rws, err := getSameVowelsWords(t)
		if err != nil {
			// エラーハンドリング
			log.Fatal(err)
		}
		res, err := json.Marshal(rws)
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(res); err != nil {
			// エラーハンドリング
			log.Fatal(err)
		}
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		log.Println(2)
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(r) != true {
			http.Error(w, "Unauthorized", 401)
			return
		}

		_, err := fmt.Fprintf(w, "Successful Basic Authentication\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/v1/roman", h)
	http.HandleFunc("/", h2)
	http.HandleFunc("/basic", h3)
	http.ListenAndServe(":8080", nil)
}
