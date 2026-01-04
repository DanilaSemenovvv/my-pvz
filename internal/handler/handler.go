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

func GetIDsInput(text string, scanner *bufio.Scanner) ([]int, error) {
	fmt.Println(text)
	scanner.Scan()

	input := strings.TrimSpace(scanner.Text())

	parts := strings.Split(input, " ")
	ids := make([]int, 0, len(parts))

	for _, p := range parts {
		p := strings.TrimSpace(p) //Удаляем лишние пробелы слева и справа от p (" 101" -> "101")

		id, err := strconv.Atoi(p) //Перевод string в int
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)

	}

	return ids, nil
}
