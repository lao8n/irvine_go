package main

import (
	"fmt"
	"sync"
)

const numPhilosophers = 5
const numChopsticks = 5

func main(){
	var wg sync.WaitGroup

	// setup chopsticks
	// use 0-4 to make modulo arithmetic work
	chopsticks := make([] *Chopstick, 0, numChopsticks)
	for i := 0; i < numChopsticks; i++ {
		chopstickChannel := make(chan bool, 1)
		// switch to id 1-5
		chopstick := &Chopstick{i + 1, chopstickChannel} 
		chopstick.available <- true
		chopsticks = append(chopsticks, chopstick) 
	}

	// setup philosophers
	// use 0-4 to make modulo arithmetic work
	philosophers := make([] *Philosopher, 0, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		// switch to id 1-5
		philosopher := &Philosopher{i + 1, chopsticks[i], chopsticks[ (i+1) % 5 ]}
		philosophers = append(philosophers, philosopher)
	}

	// setup host
	// buffer two as can process two at a time,
	var hostPermissionChannel = make(chan int, 2)
	
	// eat
	wg.Add(5)
	for i := 0; i < numPhilosophers; i++ {
		go philosophers[i].eat(&wg, hostPermissionChannel)
	}

	// have main goroutine wait for other go routines
	wg.Wait()
}

type Chopstick struct {
	id int
	available chan bool
}

type Philosopher struct {
	id int
	leftChopstick *Chopstick 
	rightChopstick *Chopstick
}

func (p *Philosopher) eat(wg *sync.WaitGroup, hostPermissionChannel chan int){
	// each philosopher eats only 3 times
	for i := 1; i <= 3; i++ {
		// permission to eat from host
		hostPermissionChannel <- p.id

		// wait for chopsticks
		select {
			// if left available first
			case <- p.leftChopstick.available:
				//fmt.Println(p.id, " has left ", p.leftChopstick.id, " waiting for right ", p.rightChopstick.id)
				// wait for right
				<- p.rightChopstick.available
			// if right available first
			case <- p.rightChopstick.available:
				//fmt.Println(p.id, " has right", p.rightChopstick.id, " waiting for left", p.leftChopstick.id)
				// wait for left
				<- p.leftChopstick.available
		}

		// eat
		fmt.Println("starting to eat  ", p.id)

		// finish eating - drop chopsticks
		fmt.Println("finishing eating ", p.id)
		p.leftChopstick.available <- true 
		p.rightChopstick.available <- true

		// release semaphore
		<- hostPermissionChannel
	}
	
	// done
	wg.Done()
}
