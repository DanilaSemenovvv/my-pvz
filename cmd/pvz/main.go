package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DanilaSemenovvv/my-pvz/internal/handler"
	"github.com/DanilaSemenovvv/my-pvz/internal/models"
	"github.com/DanilaSemenovvv/my-pvz/internal/repository"
	"github.com/DanilaSemenovvv/my-pvz/internal/service"
)

const (
	dataFile = "./internal/database/data.json"
)

func main() {
	fmt.Println(time.Now())
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("Дорогой пользователь,  выберете действие из списка представленных")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("1. Принять заказ от курьера")
	fmt.Println("2. Вернуть заказ курьеру")
	fmt.Println("3. Выдать заказ/принять возврат клиента")
	fmt.Println("4. Список заказов")
	fmt.Println("5. Список возвратов")
	fmt.Println("6. История заказов")

	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("-------------------------Ваше действие---------------------------")
	fmt.Println("-----------------------------------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		actionChoise := strings.TrimSpace(scanner.Text())
		if actionChoise == "exit" {
			fmt.Println("Досвидания")
			break
		}

		switch actionChoise {
		case "1":
			id, err := handler.GetIntInput("Введите ID-заказа: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-заказа: ", err)
				continue
			}

			userID, err := handler.GetIntInput("Введите ID-пользователя: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-пользователя: ", err)
				continue
			}
			saveDate, err := handler.GetDateInput("Введите дату хранения: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода даты хранения: ", err)
				continue
			}

			orderIssued := false //Переменная для обозначения выдан заказ или нет

			order := models.Order{
				ID:          id,
				UserID:      userID,
				SaveDate:    saveDate,
				OrderIssued: orderIssued,
			}

			err = service.ReceivingOrder(order, dataFile)
			if err != nil {
				fmt.Println(err)
			}

		case "2":
			id, err := handler.GetIntInput("Введите ID заказа для курьера", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-заказа")
				continue
			}

			err = repository.DeleteOrder(id, dataFile)
			if err != nil {
				fmt.Println(err)
			}

		case "3":
			fmt.Println(3)
		case "4":
			fmt.Println(4)
		case "5":
			fmt.Println(5)
		case "6":
			fmt.Println(6)
		default:
			fmt.Println("Неверный выбор, попробуйте снова")
		}
		fmt.Println("Выберете следующее действие или ввведите 'exit' для выхода (или нажмите сочетание клавишь ctrl+c)")
	}
}
