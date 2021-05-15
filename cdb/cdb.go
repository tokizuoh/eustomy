package cdb

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"roman"
)

type word struct {
	wordid int
	lang   string
	lemma  string
	pron   string
	pos    int
}

type romanWord struct {
	raw   string
	roman string
	// vowels string
}

// [TODO]: 関数の機能が複数（クエリ実行、構造体へのパース）なので分けたほうが良さそう
// getJapaneseWords executes a query that returns Japanese only rows.
func getJapaneseWords(db *sql.DB) ([]word, error) {
	rows, err := db.Query("select * from word where lang='jpn'")
	if err != nil {
		return []word{}, err
	}
	defer rows.Close()

	var words = []word{}
	for rows.Next() {
		var w word
		rows.Scan(&w.wordid, &w.lang, &w.lemma, &w.pron, &w.pos)
		words = append(words, w)
	}
	return words, nil
}

// debugGenerateCSV convert struct to csv.
func debugGenerateCSV(words []word) error {
	file, err := os.OpenFile("./debug_roman.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	line := []string{"raw", "roman"}
	writer.Write(line)

	for _, word := range words {
		r, err := roman.GetRomanLetters(word.lemma)
		if err != nil {
			return err
		}
		rw := romanWord{
			raw:   word.lemma,
			roman: r,
			// vowels: "TODO: convert roman to vowel",
		}
		line := []string{rw.raw, rw.roman}
		writer.Write(line)
	}
	return nil
}

func GenerateCustomDB() error {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./wnjpn.db")
	if err != nil {
		return err
	}
	// Closeしなかったらどうなる？
	defer db.Close()

	words, err := getJapaneseWords(db)
	if err != nil {
		return err
	}

	if err := debugGenerateCSV(words); err != nil {
		log.Fatal(err)
	}

	return nil
}
