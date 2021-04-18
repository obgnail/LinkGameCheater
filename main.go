package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func main22() {
	// 向上滚动：3行
	robotgo.ScrollMouse(3, `up`)
	// 向下滚动：2行
	robotgo.ScrollMouse(2, `down`)

	// 按下鼠标左键
	// 第1个参数：left(左键) / center(中键，即：滚轮) / right(右键)
	// 第2个参数：是否双击
	robotgo.MouseClick(`left`, false)

	// 按住鼠标左键
	robotgo.MouseToggle(`down`, `left`)
	// 解除按住鼠标左键
	robotgo.MouseToggle(`up`, `left`)
}
