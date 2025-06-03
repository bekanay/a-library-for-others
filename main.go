package main

import (
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
			if err != io.EOF {
				fmt.Println("Error reading line:", err)
				return
			}
		}
		
		fmt.Println(line)
		for i := 0; i < csvparser.GetNumberOfFields(); i++ {
			fmt.Println(csvparser.GetField(i))
		}
		fmt.Println()
		if err == io.EOF {
			break
		}
	}
}