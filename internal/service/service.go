package service

import (
	"fmt"
	"time"

	"github.com/DanilaSemenovvv/my-pvz/internal/models"
	"github.com/DanilaSemenovvv/my-pvz/internal/repository"
)

func ReceivingOrder(newOrder models.Order, filename string) error {
	orders, err := repository.ReadOrders(filename)
	if err != nil {
		return err
	}

	nowTime := time.Now()

	if !repository.ContainsOrders(orders, newOrder, nowTime) {
		orders = append(orders, newOrder)
		fmt.Println("Заказ принят на ПВЗ")
	} else {
		fmt.Println("Такой заказ уже был внесен в БД ПВЗ")
	}

	return repository.WriteOrder(filename, orders)

}
