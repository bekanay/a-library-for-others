package main

import (
	"errors"
	"io"
)

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
					fields, err := splitLine(string(buf))
					p.fields = fields
					if err != nil {
						return "", err
					}
					return string(buf), io.EOF
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
