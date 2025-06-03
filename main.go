package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var csvparser CSVParser = &MyCSVParser{}

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}
		
		fmt.Println(line)
		for i := 0; i < csvparser.GetNumberOfFields(); i++ {
			fmt.Println(csvparser.GetField(i))
		}
		fmt.Println()
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
	fields []string
}

func (p *MyCSVParser) ReadLine(r io.Reader) (string, error) {
	var buf []byte
	var single [1]byte

	for {
		n, err := r.Read(single[:])
		if n > 0 {
			b := single[0]
			if b == '\n' {
				break
			}
			buf = append(buf, b)
		}
		if err != nil {
			if err == io.EOF {
				if len(buf) > 0 {
					return "", io.EOF
				}
			}
			return "", err
		} 
	}

	line := string(buf)
	fields, err := splitLine(line)
	if err != nil {
		return "", err
	}
	p.fields = fields
	return line, nil
}

func (p *MyCSVParser) GetField(n int) (string, error) {
	if n < 0 || n >= len(p.fields) {
		return "", ErrFieldCount
	}
	return p.fields[n], nil
}

func (p *MyCSVParser) GetNumberOfFields() int {
	return len(p.fields)
}

func trimQuotes(word string) string {
	return word[1:len(word)-1]
}
func splitLine(line string) ([]string, error) {
	fields := []string{}
	numberOfQuotes := 0
	isQuote := false
	word := ""
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			numberOfQuotes++
			word += string(line[i])
			isQuote = !isQuote
		} else if line[i] == ',' {
			if !isQuote {
				if numberOfQuotes > 0 {
					word = trimQuotes(word)
				}
				fields = append(fields, word)
				word = ""
			} else {
				word += string(line[i])
			}
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
