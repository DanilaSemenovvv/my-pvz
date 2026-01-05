package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	repo := repository.NewFileRepository(dataFile)
	orderService, err := service.NewOrderService(repo)

	if err != nil {
		fmt.Println("Ошибка инициализации сервиса:", err)
		return
	}

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

			err = orderService.AddOrder(order)
			if err != nil {
				fmt.Println(err)
			}

		case "2":
			id, err := handler.GetIntInput("Введите ID-заказа для курьера", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода ID-заказа")
				continue
			}

			userID, err := handler.GetIntInput("Введите ID-пользователя: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода:", err)
				continue
			}

			if err := orderService.DeleteOrder(id, userID); err != nil {
				fmt.Println("Ошибка удаления:", err)
			} else {
				fmt.Println("Заказ успешно удалён и передан курьеру.")
			}

		case "3": //
			userID, err := handler.GetIntInput("Введите ID-пользователя:", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода:", err)
				continue
			}

			ids, err := handler.GetIDsInput("Введите ID заказов через пробел: ", scanner)

			if err != nil {
				fmt.Println("Ошибка ввода ID заказов:", err)
				continue
			}

			if len(ids) == 0 {
				fmt.Println("Не введено ни одного ID заказа")
				continue
			}

			fmt.Println("Выберите действие:")
			fmt.Println("1. Выдать заказы клиенту")
			fmt.Println("2. Принять возврат от клиента")
			fmt.Print("Ваш выбор (1 или 2): ")

			scanner.Scan()
			choice := strings.TrimSpace(scanner.Text())

			var action string
			switch choice {
			case "1":
				action = "issue"
			case "2":
				action = "return"
			default:
				fmt.Println("Неверный выбор действия")
				continue
			}

			success, failed, messages, err := orderService.ProcessClientOrders(userID, ids, action)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}

			fmt.Println("\n--- Результат ---")
			for _, msg := range messages {
				fmt.Println(msg)
			}

			if len(success) > 0 {
				fmt.Printf("\nУспешно обработано заказов: %d\n", len(success))
			}

			if len(failed) > 0 {
				fmt.Printf("Не удалось обработать заказов: %d\n", len(failed))
			}

		case "4":
			/*Мы получаем список заказов используя ID клиента
			и опционгальные параметры, а именно:
					1. N-последних заказов
					2. Заказы клиента которые находятся в самом ПВЗ, выданные не берутся в расчет
			*/
			userID, err := handler.GetIntInput("Введите ID-пользователя: ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода:", err)
				continue
			}
			fmt.Println("Хотите увидеть только заказы, находящиеся на ПВЗ? (y/n)")

			scanner.Scan()
			choice := strings.TrimSpace(scanner.Text())

			var onlyOnPVZ bool
			switch choice {
			case "y":
				onlyOnPVZ = true
			case "n":
				onlyOnPVZ = false

			}

			fmt.Print("Введите количество последних заказов (оставьте пустым для всех): ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())

			var limit int
			if input != "" {
				var err error
				limit, err = strconv.Atoi(input)
				if err != nil || limit < 0 {
					fmt.Println("Неверное число, показываем все заказы")
					limit = 0
				}
			} else {
				limit = 0
			}

			orderList, err := orderService.GetClientOrders(userID, onlyOnPVZ, limit)
			if err != nil {
				fmt.Println("Ошибка:", err)
			}

			if len(orderList) == 0 {
				fmt.Println("У клиента нет подходящих заказов.")
				continue
			}

			fmt.Printf("\nСписок заказов клиента #%d", userID)
			if onlyOnPVZ {
				fmt.Println(" (только на ПВЗ)")
			}
			if limit > 0 {
				fmt.Printf(" (последние %d):\n", limit)
			}

			now := time.Now()
			for i, order := range orderList {
				status := "На хранении"
				if order.OrderIssued {
					status = "Выдан клиенту"
				} else if order.SaveDate.Before(now) {
					status = "Просрочен"
				}

				fmt.Printf("%d. Заказ #%d | Хранение до: %s | %s\n",
					i+1, order.ID, order.SaveDate.Format("02.01.2006"), status)
			}

			fmt.Printf("\nВсего показано: %d заказов\n", len(orderList))

		case "5":
			now := time.Now()
			page, err := handler.GetIntInput("Введите номер страницы (1,2,...): ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода страницы:", err)
				continue
			}
			if page < 1 {
				page = 1
			}
			pageSize := 10

			returnableOrders, total, err := orderService.GetReturnableOrders(page, pageSize)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}

			if total == 0 {
				fmt.Println("Нет заказов, которые можно вернуть.")
				continue
			}

			pages := (total + pageSize - 1) / pageSize // кол-во страниц
			if len(returnableOrders) == 0 {
				fmt.Printf("Страница %d пуста. Всего заказов: %d (страниц: %d)\n", page, total, pages)
				continue
			}

			fmt.Printf("\nСписок возвратов (страница %d из %d, всего %d):\n\n", page, pages, total)

			for i, order := range returnableOrders {
				remaining := 48*time.Hour - now.Sub(order.IssuedAt)
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60

				fmt.Printf("%d. Заказ #%d | Клиент #%d | Выдан: %s | Осталось: %dч %dмин\n",
					i+1, order.ID, order.UserID, order.IssuedAt.Format("02.01.2006 15:04"), hours, minutes)
			}

		case "6":
			now := time.Now()
			page, err := handler.GetIntInput("Введите номер страницы (1,2,...): ", scanner)
			if err != nil {
				fmt.Println("Ошибка ввода страницы:", err)
				continue
			}
			if page < 1 {
				page = 1
			}
			pageSize := 10
			allOrders, total, err := orderService.GetAllOrders(page, pageSize)
			if err != nil {
				fmt.Println("Ошибка: ", err)
				continue
			}

			if total == 0 {
				fmt.Println("Нет заказов.")
				continue
			}

			pages := (total + pageSize - 1) / pageSize // кол-во страниц
			if len(allOrders) == 0 {
				fmt.Printf("Страница %d пуста. Всего заказов: %d (страниц: %d)\n", page, total, pages)
				continue
			}

			fmt.Printf("\nСписок всех заказов (страница %d из %d, всего %d):\n\n", page, pages, total)

			for i, order := range allOrders {
				status := "На хранении"
				if order.OrderIssued {
					status = "Выдан клиенту"
				} else if order.SaveDate.Before(now) {
					status = "Просрочен"
				} else if !order.IssuedAt.IsZero() { // был выдан и возвращён
					status = "Возврат принят"
				}

				fmt.Printf("%d. Заказ #%d | Клиент #%d | Хранение до: %s | Статус: %s\n",
					i+1, order.ID, order.UserID, order.SaveDate.Format("02.01.2006"), status)
			}
		default:
			fmt.Println("Неверный выбор, попробуйте снова")
		}
		fmt.Println("Выберете следующее действие или ввведите 'exit' для выхода (или нажмите сочетание клавишь ctrl+c)")
	}
}
