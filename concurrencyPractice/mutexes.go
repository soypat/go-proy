package mutexes

import (
	"sync"
)

//var y = 0

func chanIncrement(y *int,wg *sync.WaitGroup, ch chan bool) {
	ch <- true
	*y = *y + 1
	<-ch
	wg.Done()
}

// ChanDoSum Más lento que Mutex por factor de ~=3
func ChanDoSum(upTo int) int {
	y:=0
	var w sync.WaitGroup
	ch := make(chan bool, 1)
	for i := 0; i < upTo; i++ {
		w.Add(1)
		go chanIncrement(&y ,&w, ch)
	}
	w.Wait()
	return y
	//fmt.Println("final value of y", y)
	//y=0
}

//var x = 0

func mutexIncrement(x *int,wg *sync.WaitGroup, m *sync.Mutex) {
	m.Lock()
	*x = *x + 1
	m.Unlock()
	wg.Done()
}
// MutexDoSum Gotta go fast. 0.36s (Mutex) vs. 1.0s (chan)
func MutexDoSum(upTo int) int {
	x:= 0
	var w sync.WaitGroup
	var m sync.Mutex
	for i := 0; i < upTo; i++ {
		w.Add(1)
		go mutexIncrement(&x,&w, &m)

	}
	w.Wait()
	return x
	//fmt.Println("final value of x", x)
}
