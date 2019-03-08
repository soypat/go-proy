package boss

import (
"fmt"
)

type boss struct {
	firstName   string
	lastName    string
	totalLeaves int
	leavesTaken int
}

func New(firstName string, lastName string, totalLeave int, leavesTaken int) boss {
	b := boss {firstName, lastName, totalLeave, leavesTaken}
	return b
}

func (b boss) LeavesRemaining() {
	fmt.Printf("%s %s has %d leaves remaining", b.firstName, b.lastName, (b.totalLeaves - b.leavesTaken))
}
