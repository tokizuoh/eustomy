package devide

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func processing(lines [][]string, w *csv.Writer) {
	var r [][]string
	set := make(map[string]struct{})
	for _, line := range lines {
		if line[1] == "" {
			continue
		}
		if _, exist := set[line[0]]; exist {
			continue
		}
		set[line[0]] = struct{}{}
		r = append(r, line)
	}
	w.WriteAll(r)
}

func Devide() error {
	v := 1000

	inputName := "debug_roman.csv"
	file, err := os.Open(inputName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return err
	}

	outputName := "2_output_" + inputName
	file, err = os.OpenFile(outputName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	writer.Write(header)

	isEOF := false
	for {
		var lines [][]string
		for i := 0; i < v; i++ {
			line, err := reader.Read()
			if err == io.EOF {
				isEOF = true
				break
			}
			if err != nil {
				return err
			}
			lines = append(lines, line)
		}

		go processing(lines, writer)

		if isEOF {
			break
		}
	}
	return nil
}
