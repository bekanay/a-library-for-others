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
		// fmt.Println(csvparser.GetField(1))
		// fmt.Println(csvparser.GetField(2))
		// fmt.Println(csvparser.GetField(3))
		// fmt.Println(csvparser.GetNumberOfFields())
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
					line := string(buf)
					p.fields, _ = splitLine(line)
					return line, io.EOF
				}
				return "", io.EOF
			}
			return "", err
		}
	}

	line := string(buf)
	p.fields, _ = splitLine(line)
	return line, nil
}

func (p *MyCSVParser) GetField(n int) (string, error) {
	if n < 1 || n > len(p.fields) {
		return "", ErrFieldCount
	}
	return p.fields[n-1], nil
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
