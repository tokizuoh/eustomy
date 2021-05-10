package roman

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Word struct {
	Surface  string
	Furigana string
	Roman    string
}

type WordList struct {
	Word Word
}

type Result struct {
	WordList WordList
}

type ResultSet struct {
	Result Result
}

func convertStructfromXML(b []byte) (ResultSet, error) {
	rs := ResultSet{}
	if err := xml.Unmarshal(b, &rs); err != nil {
		return ResultSet{}, err
	}
	return rs, nil
}

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
	parames.Add("sentence", str)
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

	rs, err := convertStructfromXML(body)
	if err != nil {
		return "", err
	}

	return rs.Result.WordList.Word.Roman, nil
}
