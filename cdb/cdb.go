package cdb

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

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

type RomanWord struct {
	Raw    string `json:"row"`
	Roman  string `json:"roman"`
	Vowels string `json:"vowels"`
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

func extractCustomVowels(target, vowels string) string {
	res := ""
	for _, tr := range target {
		t := string([]rune{tr})
		if strings.Contains(vowels, t) {
			res += t
		}
	}
	return res
}

// debugGenerateCSV convert struct to csv.
func debugGenerateCSV(words []word) error {
	file, err := os.OpenFile("./debug_roman_vowels.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	line := []string{"raw", "roman", "vowels"}
	writer.Write(line)

	for _, word := range words {
		r, err := roman.GetRomanLetters(word.lemma)
		if err != nil {
			return err
		}
		rw := RomanWord{
			Raw:    word.lemma,
			Roman:  r,
			Vowels: extractCustomVowels(r, "aiueon"),
		}
		line := []string{rw.Raw, rw.Roman, rw.Vowels}
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

func GetSameVowelsWords(csvPath, target string) ([]RomanWord, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	targetVowel := extractCustomVowels(target, "aiueon")

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	res := []RomanWord{}
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		vowel := line[2]
		if targetVowel == vowel {
			rw := RomanWord{
				Raw:    line[0],
				Roman:  line[1],
				Vowels: line[2],
			}
			res = append(res, rw)
		}

	}

	return res, nil
}
