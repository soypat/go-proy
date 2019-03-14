//    Documentation
// https://godoc.org/github.com/nsf/termbox-go

package main

import tb "github.com/nsf/termbox-go"

func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}

	tb.Close()
}
