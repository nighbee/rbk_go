package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Process(input string) string {
	input = strings.TrimPrefix(input, "\ufeff")

	words := strings.Fields(input)
	if len(words) == 0 {
		return ""
	}

	words = handleCommands(words)

	words = handleArticles(words)

	refinedText := strings.Join(words, " ")

	refinedText = handlePunctuation(refinedText)

	refinedText = handleQuotes(refinedText)

	return strings.Join(strings.Fields(refinedText), " ")
}

func handleCommands(words []string) []string {
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		if strings.HasSuffix(word, ",") && strings.HasPrefix(word, "(") {
			cmdType := word[1 : len(word)-1]
			if i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
				nStr := strings.TrimSuffix(words[i+1], ")")
				n, err := strconv.Atoi(nStr)
				if err == nil {
					start := len(result) - n
					if start < 0 {
						start = 0
					}
					for j := start; j < len(result); j++ {
						result[j] = applyTransform(result[j], cmdType)
					}
					i++
					continue
				}
			}
		}

		if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") {
			cmd := word
			if len(result) > 0 {
				lastIdx := len(result) - 1
				applied := true
				switch cmd {
				case "(hex)":
					val, err := strconv.ParseInt(result[lastIdx], 16, 64)
					if err == nil {
						result[lastIdx] = strconv.FormatInt(val, 10)
					} else {
						applied = false
					}
				case "(bin)":
					val, err := strconv.ParseInt(result[lastIdx], 2, 64)
					if err == nil {
						result[lastIdx] = strconv.FormatInt(val, 10)
					} else {
						applied = false
					}
				case "(up)":
					result[lastIdx] = strings.ToUpper(result[lastIdx])
				case "(low)":
					result[lastIdx] = strings.ToLower(result[lastIdx])
				case "(cap)":
					result[lastIdx] = Capitalize(result[lastIdx])
				default:
					applied = false
				}
				if applied {
					continue
				}
			}
		}

		result = append(result, word)
	}
	return result
}

func handleArticles(words []string) []string {
	vowels := "aeiouhAEIOUH"
	for i := 0; i < len(words)-1; i++ {
		lowerWord := strings.ToLower(words[i])
		if lowerWord == "a" {
			nextWord := words[i+1]
			trimNext := strings.TrimLeft(nextWord, "'")
			if len(trimNext) > 0 && strings.ContainsRune(vowels, rune(trimNext[0])) {
				if words[i] == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}
	return words
}

func handlePunctuation(text string) string {
	reClose := regexp.MustCompile(`\s+([.,!?:;]+)`)
	text = reClose.ReplaceAllString(text, "$1")

	reOpen := regexp.MustCompile(`([.,!?:;]+)([^\s.,!?:;])`)
	text = reOpen.ReplaceAllString(text, "$1 $2")

	return text
}

func handleQuotes(text string) string {
	var result strings.Builder
	isOpening := true

	parts := strings.Split(text, "'")
	for i, part := range parts {
		if i == 0 {
			result.WriteString(strings.TrimRight(part, " "))
		} else {
			if isOpening {
				result.WriteString(" '")
				result.WriteString(strings.TrimSpace(part))
				isOpening = false
			} else {
				result.WriteString("'")
				remaining := strings.TrimLeft(part, " ")
				if remaining != "" {
					result.WriteString(" ")
				}
				result.WriteString(remaining)
				isOpening = true
			}
		}
	}

	return result.String()
}

func applyTransform(s, cmd string) string {
	switch cmd {
	case "up":
		return strings.ToUpper(s)
	case "low":
		return strings.ToLower(s)
	case "cap":
		return Capitalize(s)
	default:
		return s
	}
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	first := unicode.ToUpper(r[0])
	rest := strings.ToLower(string(r[1:]))
	return string(first) + rest
}
