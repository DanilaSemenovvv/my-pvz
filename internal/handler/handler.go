package handler

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "02.01.2006"
)

func GetIntInput(text string, scanner *bufio.Scanner) (int, error) {
	fmt.Println(text)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	return strconv.Atoi(input)
}

func GetDateInput(text string, scanner *bufio.Scanner) (time.Time, error) {
	fmt.Println(text)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	return time.Parse(dateFormat, input)
}
