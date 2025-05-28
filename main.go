package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var csvparser CSVParser = &MyCSVParser{
		reader: bufio.NewReader(file),
	}
	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					fmt.Println(line)
				}
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}
		fmt.Println(line)
		fmt.Println(csvparser.GetField(1))
		fmt.Println(csvparser.GetField(2))
		fmt.Println(csvparser.GetField(3))
		fmt.Println(csvparser.GetNumberOfFields())
	}
}

type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

type MyCSVParser struct {
	reader *bufio.Reader
	fields []string
}

func (p *MyCSVParser) ReadLine(r io.Reader) (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF && len(line) > 0 {
			line = strings.TrimSpace(line)
			return line, io.EOF
		}
		return "", err
	}
	line = strings.TrimSpace(line)
	p.fields, err = splitLine(line)
	if err != nil {
		return "", err
	}
	return line, nil
}

func (p *MyCSVParser) GetField(n int) (string, error) {
	if n < 0 || n > len(p.fields) {
		return "", ErrFieldCount
	}
	field := p.fields[n-1]
	return field, nil
}

func (p *MyCSVParser) GetNumberOfFields() int {
	return len(p.fields)
}

func splitLine(line string) ([]string, error) {
	fields := []string{}
	isQuote := false
	word := ""
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			if !isQuote {
				isQuote = true
			} else {
				fields = append(fields, word)
			}
		} else if line[i] == ',' && !isQuote {
			fields = append(fields, word)
			word = ""
		} else {
			word += string(line[i])
		}
	}
	if isQuote {
		return nil, ErrQuote
	}
	fields = append(fields, word)
	return fields, nil
}
