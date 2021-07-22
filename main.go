package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var equationRegex = regexp.MustCompile(`[0-9]+[\+\-\*\:][0-9]+[=][?]`)
var operationRegex = regexp.MustCompile(`[0-9]+[\+\-\*\:][0-9]+`)
var operators = []string{"+", "-", "*", ":"}

func main() {
	commandReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the name of an input file (including extension): ")
	inputFileName, _, err := commandReader.ReadLine()
	if err != nil {
		fmt.Println(err)
		return
	}
	inputFile, err := os.OpenFile(string(inputFileName), os.O_RDONLY, 0777)

	defer inputFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Enter the name of an output file (including extension): ")
	outputFileName, _, err := commandReader.ReadLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	outputFile, err := os.OpenFile(string(outputFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer outputFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		equasion := equationRegex.FindString(string(line))
		if equasion == "" {
			continue
		}

		operation := operationRegex.FindString(equasion)
		if operation == "" {
			continue
		}

		numbers := []string{}
		var currentOperator string

		for _, op := range operators {
			numbers = strings.Split(operation, op)
			if len(numbers) > 1 {
				currentOperator = op
				break
			}
		}

		if currentOperator == "" || len(numbers) < 2 {
			continue
		}

		var answer float32

		first, err := strconv.Atoi(numbers[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		second, err := strconv.Atoi(numbers[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		switch currentOperator {
		case "+":
			answer = float32(first + second)
		case "-":
			answer = float32(first - second)
		case "*":
			answer = float32(first * second)
		case ":":
			answer = float32(first / second)
		default:
			fmt.Println("Unsupported operator in line:\n", equasion)
			return
		}

		outputString := fmt.Sprintf("%d%v%d=%0.f", first, currentOperator, second, answer)

		_, err = writer.WriteString(outputString + "\n")

		if err != nil {
			fmt.Println(err)
			break
		}
	}

	writer.Flush()
	fmt.Println("Results appended to:", string(outputFileName))
}
