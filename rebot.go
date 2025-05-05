package main

import (
	"fmt"
	"time"
)

type Robot struct {
	ID         int
	Controller *Controller
	CancelChan chan bool
	IsBusy     bool
	CurrentJob *Order
}

func (r *Robot) Start() {
	go func() {
		for {
			order := r.Controller.GetNextOrderForRobot(r.ID)
			if order == nil {
				// 没有订单，等待新订单
				time.Sleep(1 * time.Second)
				continue
			}
			r.IsBusy = true
			r.CurrentJob = order
			logChannel <- fmt.Sprintf("[Robot %d] 开始处理订单 %d (VIP: %v)\n", r.ID, order.ID, order.IsVIP)

			select {
			case <-time.After(10 * time.Second):
				// 完成订单
				r.Controller.MarkOrderComplete(*order)
				logChannel <- fmt.Sprintf("[Robot %d] 完成订单 %d\n", r.ID, order.ID)
				r.IsBusy = false
				r.CurrentJob = nil
			case <-r.CancelChan:
				// 停止处理
				logChannel <- fmt.Sprintf("[Robot %d] 停止处理订单 %d\n", r.ID, order.ID)
				r.Controller.ReAddPendingOrder(*order)
				r.IsBusy = false
				r.CurrentJob = nil
				return
			}
		}
	}()
}
