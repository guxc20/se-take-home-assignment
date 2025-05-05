package main

import (
	"fmt"
	"sync"
)

type Controller struct {
	PendingQueue  *OrderQueue
	CompleteQueue []Order
	Robots        []*Robot
	orderCounter  int
	mu            sync.Mutex
}

func NewController() *Controller {
	return &Controller{
		PendingQueue:  NewOrderQueue(),
		CompleteQueue: []Order{},
		Robots:        []*Robot{},
		orderCounter:  0,
	}
}

func (c *Controller) CreateOrder(isVIP bool) {
	c.mu.Lock()
	c.orderCounter++
	order := Order{
		ID:    c.orderCounter,
		IsVIP: isVIP,
	}
	c.PendingQueue.AddOrder(order)
	c.mu.Unlock()
	logChannel <- fmt.Sprintf("[系统] 新订单 %d 已加入 (VIP: %v)\n", order.ID, order.IsVIP)
}

func (c *Controller) AddRobot() {
	c.mu.Lock()
	defer c.mu.Unlock()

	robot := &Robot{
		ID:         len(c.Robots) + 1,
		Controller: c,
		CancelChan: make(chan bool),
	}
	c.Robots = append(c.Robots, robot)
	logChannel <- fmt.Sprintf("[系统] 添加机器人 %d\n", robot.ID)
	robot.Start()
}

func (c *Controller) RemoveRobot() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Robots) == 0 {
		logChannel <- fmt.Sprintf("[系统] 没有机器人可以移除\n")
		return
	}
	robot := c.Robots[len(c.Robots)-1]
	c.Robots = c.Robots[:len(c.Robots)-1]

	if robot.IsBusy {
		robot.CancelChan <- true
	} else {
		close(robot.CancelChan)
	}
	logChannel <- fmt.Sprintf("[系统] 移除机器人 %d\n", robot.ID)
}

func (c *Controller) GetNextOrderForRobot(robotID int) *Order {
	return c.PendingQueue.PopOrder()
}

func (c *Controller) MarkOrderComplete(order Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompleteQueue = append(c.CompleteQueue, order)
}

func (c *Controller) ReAddPendingOrder(order Order) {
	c.PendingQueue.ReAddOrder(order)
}

func (c *Controller) PrintStatus() {
	c.mu.Lock()
	defer c.mu.Unlock()

	logChannel <- fmt.Sprintf("====== 当前状态 ======")
	logChannel <- fmt.Sprintf("待处理订单:")
	for _, o := range c.PendingQueue.GetOrders() {
		logChannel <- fmt.Sprintf("- 订单 %d (VIP: %v)", o.ID, o.IsVIP)
	}
	logChannel <- fmt.Sprintf("已完成订单:")
	for _, o := range c.CompleteQueue {
		logChannel <- fmt.Sprintf("- 订单 %d (VIP: %v)", o.ID, o.IsVIP)
	}
	logChannel <- fmt.Sprintf("机器人状态:")
	for _, r := range c.Robots {
		if r.IsBusy {
			logChannel <- fmt.Sprintf("- 机器人 %d: 正在处理订单 %d", r.ID, r.CurrentJob.ID)
		} else {
			logChannel <- fmt.Sprintf("- 机器人 %d: 空闲", r.ID)
		}
	}
	logChannel <- fmt.Sprintf("====================\n")
}
