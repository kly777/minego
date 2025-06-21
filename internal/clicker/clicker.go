// clicker.go
package clicker

import (
	"image"
	"time"

	"minego/pkg/winapi/click"
)

// ClickType 点击类型枚举
type ClickType int

const (
	LeftClick ClickType = iota
	RightClick
	SpecialClick
)

// ClickTask 点击任务结构体
type ClickTask struct {
	Point image.Point
	Type  ClickType
	Delay time.Duration
}

// Clicker 点击处理器
type Clicker struct {
	taskChan chan ClickTask
	done     chan struct{}
}

// NewClicker 创建点击处理器
func NewClicker(bufferSize int) *Clicker {
	return &Clicker{
		taskChan: make(chan ClickTask, bufferSize),
		done:     make(chan struct{}),
	}
}

// Start 启动点击协程
func (c *Clicker) Start() {
	go func() {
		defer close(c.done)
		for task := range c.taskChan {
			switch task.Type {
			case LeftClick:
				click.Click(task.Point)
			case RightClick:
				click.RightClick(task.Point)
			case SpecialClick:
				click.Click(task.Point) // 特殊点击逻辑可扩展
			}

			if task.Delay > 0 {
				time.Sleep(task.Delay)
			}
		}
	}()
}

// Stop 停止点击协程
func (c *Clicker) Stop() {
	close(c.taskChan)
	<-c.done
}

// SendTasks 发送批量点击任务
func (c *Clicker) SendTasks(points []image.Point, clickType ClickType, delay time.Duration) {
	for _, point := range points {
		c.taskChan <- ClickTask{
			Point: point,
			Type:  clickType,
			Delay: delay,
		}
	}
}
