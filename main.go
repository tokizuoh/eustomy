package main

import (
	"cdb"
	"encoding/json"
	"log"
	"net/http"
)

func getSameVowelsWords(s string) ([]cdb.RomanWord, error) {
	rws, err := cdb.GetSameVowelsWords("debug_roman_vowels.csv", s)
	if err != nil {
		return nil, err
	}

	return rws, nil
}

func main() {
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
	http.HandleFunc("/v1/roman", h)
	http.HandleFunc("/", h2)
	http.ListenAndServe(":8080", nil)
}
