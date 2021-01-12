package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilname := flag.String("csv", "problems.csv", " a csv file in the format of question answer")

	timeLimit := flag.Int("limit", 30, " the time limit for the quiss in second")

	flag.Parse()

	file, err := os.Open(*csvFilname)
	if err != nil {
		exit(fmt.Sprintf("Faile to open the csv file: %s\n", *csvFilname))
		os.Exit(1)
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		exit("Faile to parse the provider csv file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	fmt.Println(problems)

	correct := 0

problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d:%s=\n", i+1, problem.question)

		//create a channel
		answerCh := make(chan string)

		//go routine anonymous function
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
