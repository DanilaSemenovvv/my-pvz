package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/DanilaSemenovvv/my-pvz/internal/models"
)

const (
	invalidIndex    = -1
	filePermissions = 0644
)

func ReadOrders(file string) ([]models.Order, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Order{}, nil
		}
		return nil, err
	}

	var orders []models.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func ContainsOrders(orders []models.Order, newOrder models.Order, now time.Time) bool {
	for _, order := range orders {
		if order.ID == newOrder.ID && order.UserID == newOrder.UserID { //исправить логику проверки на прохождение заказ
			if order.SaveDate.After(now) {
				return true
			}
		}
	}

	return false
}

func WriteOrder(filename string, orders []models.Order) error {
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, filePermissions)
}

func StatusCheck(orders []models.Order, id int, now time.Time) (bool, int, error) {
	for i, order := range orders { //в чем разница между конструкцией for _, order = и for order =
		if order.ID == id {
			if order.OrderIssued {
				return false, invalidIndex, fmt.Errorf("заказ с ID %d уже выдан", id)
			}
			if order.SaveDate.Before(now) {
				return false, invalidIndex, fmt.Errorf("заказ с ID %d еще не просрочен", id)
			}

			return true, i, nil
		}
	}

	return false, invalidIndex, fmt.Errorf("заказ с ID %d не найден или уже выдан", id)
}

func DeleteOrder(id int, filename string) error {
	orders, err := ReadOrders(filename)
	if err != nil {
		return err
	}

	nowTime := time.Now()

	status, index, err := StatusCheck(orders, id, nowTime)
	if err != nil {
		return err
	}

	if status {
		orders = append(orders[:index], orders[index+1:]...)
	} else {
		return fmt.Errorf("невозможно удалить заказ с ID %d: не пройдены проверки", id)
	}

	return WriteOrder(filename, orders)
}
