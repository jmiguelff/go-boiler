package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	maxEating     = 2
	fullnessLimit = 3
)

// chopsticks is an object structure that serves as a tool to eat
type chopsticks struct {
	sync.Mutex
	ID     int
	isUsed bool
}

// newChopsticks creates a Chopsticks object with the given id label.
func newChopsticks(id int) *chopsticks {
	return &chopsticks{
		ID:     id,
		isUsed: false,
	}
}

// use is to set the Chopsticks status to "In-Use"
func (c *chopsticks) use() {
	c.Lock()
	defer c.Unlock()
	c.isUsed = true
}

// free is to set the Chopsticks status to "Free"
func (c *chopsticks) free() {
	c.Lock()
	defer c.Unlock()
	c.isUsed = false
}

// isInUse is to check the Chopsticks status, currently "In-Use" or "Free".
func (c *chopsticks) isInUse() bool {
	c.Lock()
	defer c.Unlock()
	return c.isUsed
}

// host is the structure serving foods and chopsticks
type host struct {
	sync.Mutex
	isEating int
	chops    []*chopsticks
}

// newHost creates the host object with person quantity.
func newHost(quantity int) *host {
	h := &host{
		isEating: 0,
		chops:    []*chopsticks{},
	}
	for i := 0; i < int(quantity); i++ {
		h.chops = append(h.chops, newChopsticks(i))
	}

	return h
}

// requestChopsticks is to allows a customer to eat using an available
// chopstick.
func (h *host) requestChopsticks() *chopsticks {
	h.Lock()
	defer h.Unlock()

	// table is full host don't let anyone else in
	if h.isEating >= maxEating {
		return nil
	}

	// permit to eat. Scan for available chopsticks
	c := h.seekChopsticks()
	c.use()
	h.isEating++
	return c
}

// seekChopsticks look for available chopsticks
func (h *host) seekChopsticks() *chopsticks {
	for i := range h.chops {
		if !h.chops[i].isInUse() && !h.chops[(i+1)%5].isInUse() {
			return h.chops[i]
		}
	}
	return nil
}

// returnChopsticks is to allow a customer to place back chopsticks when
// he/she is done eating
func (h *host) returnChopsticks(c *chopsticks) {
	h.Lock()
	defer h.Unlock()
	h.isEating--
	c.free()
}

func eat(id int, fullness chan int, h *host) {

	eatCount := fullnessLimit

	for {
		chops := h.requestChopsticks()
		if chops == nil {
			continue
		}

		fmt.Printf("starting to eat <%d> <count %d/%d>\n", id, fullnessLimit-eatCount, fullnessLimit)

		time.Sleep(time.Second)

		eatCount--

		fmt.Printf("finishing eating <%d> <count %d/%d>\n", id, fullnessLimit-eatCount, fullnessLimit)

		h.returnChopsticks(chops)

		if eatCount == 0 {
			fullness <- id
			return
		}
	}
}

func main() {

	nbrOfPhi := 5

	h := newHost(nbrOfPhi)
	c1 := make(chan int)
	fullness := 0

	for i := 0; i < nbrOfPhi; i++ {
		go eat(i, c1, h)
	}

	// Wait for fullness
	for {
		select {
		case person := <-c1:
			fullness++
			fmt.Printf("Philosopher %d is full\n", person)
			if fullness == nbrOfPhi {
				fmt.Printf("All are full.\n[ ENDED ]\n")
				return
			}
		}
	}
}
