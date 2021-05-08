package romen

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// TODO: 動的にする
const TARGET = "天気"

// getRomanLetters convert Kanji , Hiragana and Katakana to RomanLetters.
func GetRomanLetters(str string) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	baseURL := "https://jlp.yahooapis.jp/FuriganaService/V1/furigana"
	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return "", nil
	}

	parames := request.URL.Query()

	appID := os.Getenv("APP_ID")
	parames.Add("appid", appID)
	parames.Add("grade", "1")
	parames.Add("sentence", TARGET)
	request.URL.RawQuery = parames.Encode()

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// TODO: XML -> structに変換する
	log.Println(string(body))

	return "", nil
}
