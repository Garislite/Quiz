package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	score := 0
	total := 0

	for {
		fmt.Println("\nQUIZ GAME")
		fmt.Println("=========")
		if total > 0 {
			fmt.Printf("Score: %d/%d\n", score, total)
		}
		fmt.Println("1. New question")
		fmt.Println("2. Exit")
		fmt.Print("> ")

		var choice string
		fmt.Scan(&choice)

		if choice == "2" {
			fmt.Printf("\nFinal score: %d/%d\n", score, total)
			fmt.Println("Goodbye!")
			break
		}

		if choice != "1" {
			fmt.Println("Invalid choice")
			continue
		}

		resp, err := http.Get("https://opentdb.com/api.php?amount=1&type=multiple")
		if err != nil {
			fmt.Println("Error getting question")
			continue
		}

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if err != nil {
			fmt.Println("Error parsing response")
			continue
		}

		results, ok := result["results"].([]interface{})
		if !ok || len(results) == 0 {
			fmt.Println("No questions available")
			continue
		}

		q, ok := results[0].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid question format")
			continue
		}

		question, _ := q["question"].(string)
		correct, _ := q["correct_answer"].(string)
		incorrect, _ := q["incorrect_answers"].([]interface{})

		clean := func(s string) string {
			s = strings.ReplaceAll(s, "&quot;", "\"")
			s = strings.ReplaceAll(s, "&#039;", "'")
			s = strings.ReplaceAll(s, "&amp;", "&")
			s = strings.ReplaceAll(s, "&eacute;", "e")
			s = strings.ReplaceAll(s, "&ldquo;", "\"")
			s = strings.ReplaceAll(s, "&rdquo;", "\"")
			return s
		}

		question = clean(question)
		correct = clean(correct)

		answers := []string{correct}
		for _, ans := range incorrect {
			answers = append(answers, clean(ans.(string)))
		}

		// Простое перемешивание
		if len(answers) > 1 {
			answers[0], answers[1] = answers[1], answers[0]
		}

		fmt.Printf("\n%s\n\n", question)

		for i, ans := range answers {
			fmt.Printf("%d. %s\n", i+1, ans)
		}

		fmt.Print("\nYour answer (1-4): ")
		var userChoice int
		_, err = fmt.Scan(&userChoice)

		if err != nil || userChoice < 1 || userChoice > 4 {
			fmt.Println("Invalid input")
			continue
		}

		total++
		if answers[userChoice-1] == correct {
			fmt.Println("Correct!")
			score++
		} else {
			fmt.Printf("Wrong! Answer: %s\n", correct)
		}
	}
}
