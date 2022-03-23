package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	values := parseArguments()

	f := openFile(values)
	defer f.Close()

	records := readFile(f)

	correct, count := calcResults(records)

	fmt.Println("You got", correct, "corrected answers, out of", count, "questions")

}

func parseArguments() []string {
	flag.Parse()
	values := flag.Args()

	if len(values) > 2 {
		fmt.Println("Usage: ./quizgame <filename> <timer>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if len(values) == 0 {
		values = append(values, "problems.csv", "30")
	}
	return values

}

func openFile(values []string) *os.File {
	f, err := os.Open(values[0])

	if err != nil {
		fmt.Println("error opening file: err:", err)
		os.Exit(1)
	}
	return f

}

func readFile(f *os.File) [][]string {
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records

}

func calcResults(records [][]string) (int, int) {
	correct, count := 0, 0
	values := parseArguments()
	timeLimit, _ := strconv.Atoi(values[1])

	fmt.Println("Press Enter when you're ready to start")

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	for _, record := range records {
		fmt.Println("Question: what's the result of", record[0]+"?")
		answerCh := make(chan string)

		go func() {
			count++
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
			}
			input = strings.TrimSuffix(input, "\n")
			answerCh <- input
		}()
		select {
		case <-timer.C:
			return correct, count
		case input := <-answerCh:
			if input == record[1] {
				correct++
			}
		}

	}
	return correct, count
}
