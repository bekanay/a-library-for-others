package main

func trimQuotes(word string) string {
	return word[1 : len(word)-1]
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
				if numberOfQuotes > 0 && word[0] == '"' && word[len(word)-1] == '"'{
					word = trimQuotes(word)
					numberOfQuotes = 0
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
	if numberOfQuotes > 0 && word[0] == '"' && word[len(word)-1] == '"'{
		word = trimQuotes(word)
		numberOfQuotes = 0
	}
	fields = append(fields, word)
	return fields, nil
}
