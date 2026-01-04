package service

import (
	"fmt"
	"time"

	"github.com/DanilaSemenovvv/my-pvz/internal/models"
	"github.com/DanilaSemenovvv/my-pvz/internal/repository"
)

const (
	ActionIssue  = "issue"
	ActionReturn = "return"
)

type OrderService struct {
	repo         repository.Repository
	orders       []models.Order   // данные для сохранения в JSON
	index        map[OrderKey]int // быстрый поиск по id и userID
	byUserID     map[int][]int    //хранение по UserID  списка заказов
	issuedOrders []int            //храним индексы выданных заказов
}

type OrderKey struct {
	ID     int
	UserID int
}

func NewOrderService(repo repository.Repository) (*OrderService, error) {
	orders, err := repo.FindAll()
	if err != nil {
		return nil, err
	}

	s := &OrderService{
		repo:     repo,
		orders:   orders,
		index:    make(map[OrderKey]int, len(orders)),
		byUserID: make(map[int][]int),
	}

	for i, o := range orders {
		key := OrderKey{ID: o.ID, UserID: o.UserID}
		s.index[key] = i

		s.byUserID[o.UserID] = append(s.byUserID[o.UserID], i)
	}

	for i, o := range orders {
		if o.OrderIssued {
			s.issuedOrders = append([]int{i}, s.issuedOrders...)
		}
	}

	return s, nil
}

func (s *OrderService) AddOrder(o models.Order) error {
	key := OrderKey{ID: o.ID, UserID: o.UserID}
	now := time.Now()

	if idx, ok := s.index[key]; ok {
		existing := s.orders[idx]
		if existing.SaveDate.After(now) {
			return fmt.Errorf("заказ уже существует и ещё не просрочен")
		}
	}

	s.orders = append(s.orders, o)
	newIdx := len(s.orders) - 1
	s.byUserID[o.UserID] = append(s.byUserID[o.UserID], newIdx)
	s.index[key] = len(s.orders) - 1
	return s.repo.SaveAll(s.orders)
}

func (s *OrderService) DeleteOrder(id int, userID int) error {
	key := OrderKey{ID: id, UserID: userID}
	now := time.Now()

	idx, ok := s.index[key]
	if !ok {
		return fmt.Errorf("Заказ не найден")
	}

	order := s.orders[idx]

	if order.OrderIssued {
		return fmt.Errorf("Заказ уже выдан")
	}

	if order.SaveDate.Before(now) {
		return fmt.Errorf("Заказ еще не просрочен")
	}

	lastIndex := len(s.orders) - 1

	if idx != lastIndex {
		lastOrder := s.orders[lastIndex]
		s.orders[idx] = lastOrder

		lastKey := OrderKey{ID: lastOrder.ID, UserID: lastOrder.UserID}
		s.index[lastKey] = idx
	}

	s.orders = s.orders[:lastIndex]

	delete(s.index, key)
	s.rebuildByUserID()

	return s.repo.SaveAll(s.orders)
}

func (s *OrderService) rebuildByUserID() {
	s.byUserID = make(map[int][]int)
	for i, o := range s.orders {
		s.byUserID[o.UserID] = append(s.byUserID[o.UserID], i)
	}
}

func (s *OrderService) ProcessClientOrders(userID int, ordersIDs []int, action string) ([]int, []int, []string, error) {
	now := time.Now()

	success := []int{}
	failed := []int{}
	messages := []string{}

	for _, id := range ordersIDs {
		key := OrderKey{ID: id, UserID: userID}
		idx, exist := s.index[key]
		if !exist {
			failed = append(failed, id)
			messages = append(messages, fmt.Sprintf("Заказ %d: не найден или принадлежит другому клиенту", id))
			continue
		}

		order := &s.orders[idx] // берём указатель, чтобы менять оригинал

		if action == ActionIssue {
			// === Правила выдачи ===
			if order.OrderIssued {
				failed = append(failed, id)
				messages = append(messages, fmt.Sprintf("Заказ %d: уже выдан", id))
				continue
			}

			if order.SaveDate.Before(now) {
				failed = append(failed, id)
				messages = append(messages, fmt.Sprintf("Заказ %d: истёк срок хранения", id))
				continue
			}

			s.issuedOrders = append([]int{idx}, s.issuedOrders...)
			order.OrderIssued = true
			order.IssuedAt = now
			success = append(success, id)
			messages = append(messages, fmt.Sprintf("Заказ %d: успешно выдан клиенту", id))

		} else if action == ActionReturn {
			// === Правила возврата ===
			if !order.OrderIssued {
				failed = append(failed, id)
				messages = append(messages, fmt.Sprintf("Заказ %d: ещё не был выдан (нельзя принять возврат)", id))
				continue
			}

			if now.Sub(order.IssuedAt) > 48*time.Hour {
				failed = append(failed, id)
				messages = append(messages, fmt.Sprintf("Заказ %d: прошло более 2 суток с момента выдачи", id))
				continue
			}

			order.OrderIssued = false
			order.IssuedAt = time.Time{}
			success = append(success, id)
			messages = append(messages, fmt.Sprintf("Заказ %d: возврат принят", id))

		} else {
			return nil, nil, nil, fmt.Errorf("неизвестное действие: %s", action)
		}
	}

	if len(success) > 0 {
		if err := s.repo.SaveAll(s.orders); err != nil {
			return nil, nil, nil, fmt.Errorf("ошибка сохранения данных: %v", err)
		}
	}

	return success, failed, messages, nil
}

func (s *OrderService) GetClientOrders(userID int, onlyOnPVZ bool, limit int) ([]models.Order, error) {
	indices, ok := s.byUserID[userID]
	if !ok {
		return nil, nil
	}

	cap := len(indices)
	if limit > 0 && limit < cap {
		cap = limit
	}

	result := make([]models.Order, 0, cap)

	for i := len(indices) - 1; i >= 0; i-- {
		idx := indices[i]
		order := s.orders[idx]

		if onlyOnPVZ && order.OrderIssued {
			continue
		}

		result = append(result, order)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (s *OrderService) GetReturnableOrders(page, pageSize int) ([]models.Order, int, error) {
	now := time.Now()
	threshold := now.Add(-48 * time.Hour)

	returnable := make([]models.Order, 0)

	for _, idx := range s.issuedOrders {
		order := s.orders[idx]
		if order.IssuedAt.Before(threshold) {
			break
		}
		returnable = append(returnable, order)
	}

	total := len(returnable)

	start := (page - 1) * pageSize
	if start >= total {
		return nil, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	return returnable[start:end], total, nil
}
