package race

import "fmt"

// import "time"
import "sync"

// Data race possible
func Race(i int) {
	var m sync.Mutex
	x := 1

	go func() {
		m.Lock()
		x = 2
		m.Unlock()
	}()

	//	time.Sleep(time.Second)
	x = 3
	m.Lock()
	fmt.Printf("\n%d", x)
	m.Unlock()

}

func main5() {

	for i := 0; i < 10; i++ {
		Race(i)
	}

}
