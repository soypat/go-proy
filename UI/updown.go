package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
)

var keepGoing bool = true
var escChar string = "x"
var lastpressed rune = 'x'
var charByte []byte
func main() {
	fmt.Println("Start program")
	for keepGoing {
		char,err := getChar()
		if lastpressed==char && lastpressed!='x'{
			charByte=[]byte(string(char))
			fmt.Println("Same key has been pressed:",char,"\nByte format:",charByte)
		}
		//fmt.Println("Ascii: ",ascii,"\nKeycode: ",keyCode)
			if err!=nil {
			fmt.Println("Error")
			err=nil
		}
		charstring:= string(char)
		if charstring==escChar{
			keepGoing=false
		}
		lastpressed=char
	}

}

func getChar() (rune, error) {
	char, _, err := keyboard.GetSingleKey()
	if (err != nil) {
		return char,err
	}
	fmt.Printf("You pressed: %q\r\n", char)
	return char,nil
}



/*
faster!
func runesToUTF8Manual2(rs []rune) []byte {
    size := 0
    for _, r := range rs {
        size += utf8.RuneLen(r)
    }

    bs := make([]byte, size)

    count := 0
    for _, r := range rs {
        count += utf8.EncodeRune(bs[count:], r)
    }

    return bs
}


 */