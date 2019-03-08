package main

import (
	"fmt"
	"html"
	"strconv"
)
func SprintEmoji(charIndex *int) string {
	return html.UnescapeString("&#" + strconv.Itoa(*charIndex) + ";")
	//emoji := [][]int{
	//	// Emoticons icons.
	//	{128513, 128591},
	//	// Dingbats.
	//	{9986, 10160},
	//	// Transport and map symbols.
	//	{128640, 128704},
	//}
}
func main() {
	// Hexadecimal ranges from: http://apps.timwhitlock.info/emoji/tables/unicode
	emoji := [][]int{
		// Emoticons icons.
		{128513, 128591},
		// Dingbats.
		{9986, 10160},
		// Transport and map symbols.
		{128640, 128704},
	}

	for _, value := range emoji {
		for x := value[0]; x < value[1]; x++ {
			// Unescape the string (HTML Entity -> String).
			// Display the emoji.
			fmt.Println(SprintEmoji(&x))
		}
	}
}
