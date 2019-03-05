package mutexes

import "time"
import "fmt"

func ping1(c chan<- string) {
	t:=time.NewTicker(time.Millisecond*1000)
	for {
		c <- "ping on channel 1"
		<- t.C
	}
}

func ping2(c chan<- string) {
	t:=time.NewTicker(time.Millisecond*579)
	for {
		c <- "ping on channel 2"
		<- t.C
	}
}

func stopper(stop chan bool) {
	time.Sleep(time.Second*10)
	stop <- true
}

func main() {
	stop := make(chan bool)
	channel1 := make(chan string)
	channel2 := make(chan string)
	go stopper(stop)
	go ping1(channel1)
	go ping2(channel2)
	for {
		select {
		case <-stop:
			fmt.Println("STOOOOOOP!!!!!!!")
			return
		case msg1 := <-channel1:
			fmt.Println("received", msg1)
		case msg2 := <-channel2:
			fmt.Println("received", msg2)
		}
	}

}