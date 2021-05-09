package cdb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type word struct {
	wordid int
	lang   string
	lemma  string
	pron   string
	pos    int
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

	for _, word := range words {
		log.Println(word)
	}

	return nil
}
