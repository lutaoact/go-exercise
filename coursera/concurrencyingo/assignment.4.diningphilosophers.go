package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ChopS struct{ sync.Mutex }

type Philo struct {
	num             int // each philosopher is numbered, 1 through 5
	leftCS, rightCS *ChopS
}

func (p Philo) eat(wg *sync.WaitGroup, host chan struct{}) {
	defer wg.Done()
	for i := 0; i < 3; i++ {
		host <- struct{}{}        // get permission from host
		if rand.Float64() < 0.5 { // pick up the chopsticks in any order
			p.leftCS.Lock()
			p.rightCS.Lock()
		} else {
			p.rightCS.Lock()
			p.leftCS.Lock()
		}

		fmt.Printf("starting to eat %d\n", p.num)
		time.Sleep(1 * time.Second) // delay 1 second for having dinner
		fmt.Printf("finishing eating %d\n", p.num)

		p.rightCS.Unlock()
		p.leftCS.Unlock()

		<-host
	}
}

func main() {
	CSticks := make([]*ChopS, 5)
	for i := 0; i < 5; i++ {
		CSticks[i] = new(ChopS)
	}

	philos := make([]*Philo, 5)
	for i := 0; i < 5; i++ {
		philos[i] = &Philo{i + 1, CSticks[i], CSticks[(i+1)%5]}
	}

	host := make(chan struct{}, 2)
	wg := &sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go philos[i].eat(wg, host)
	}
	wg.Wait()
}
