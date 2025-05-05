package main

import "sync"

type Order struct {
	ID    int
	IsVIP bool
}

type OrderQueue struct {
	Orders []Order
	mu     sync.Mutex
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		Orders: []Order{},
	}
}

func (q *OrderQueue) AddOrder(order Order) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if order.IsVIP {
		// 插入到最后一个 VIP 订单之后
		pos := 0
		for i, o := range q.Orders {
			if o.IsVIP {
				pos = i + 1
			}
		}
		// 插入到 pos
		q.Orders = append(q.Orders[:pos], append([]Order{order}, q.Orders[pos:]...)...)
	} else {
		q.Orders = append(q.Orders, order)
	}
}

func (q *OrderQueue) PopOrder() *Order {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Orders) == 0 {
		return nil
	}
	order := q.Orders[0]
	q.Orders = q.Orders[1:]
	return &order
}

func (q *OrderQueue) ReAddOrder(order Order) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 回到队首
	q.Orders = append([]Order{order}, q.Orders...)
}

func (q *OrderQueue) GetOrders() []Order {
	q.mu.Lock()
	defer q.mu.Unlock()

	return append([]Order{}, q.Orders...) // copy
}
