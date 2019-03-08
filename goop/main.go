package main

import (
	"boss"
	"employee"
)

func main() {
	e := employee.New("Sam", "Adolf", 30, 20)
	e.LeavesRemaining()
	b := boss.New("Dorp","Zerp",20,23)
	b.LeavesRemaining()
}