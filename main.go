package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var balance int
var transactionNo int
var mutx sync.Mutex

func main() {
	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	balanceChan := make(chan int)
	tranChan := make(chan bool)

	balance = 1000
	transactionNo = 0
	fmt.Println("Staring balance: $", balance)

	wg.Add(1)
	for i := 0; i < 100; i++ {

		go func(ii int) {
			transactionAmount := rand.Intn(25)
			balanceChan <- transactionAmount

			if ii == 99 {
				fmt.Println("Should be quitting time")
				tranChan <- true
				close(balanceChan)
				wg.Done()
			}
		}(i)
	}
	// go transaction(0)

	breakPoint := false

	for {

		if breakPoint == true {
			break
		}

		select {

		case amt := <-balanceChan:
			fmt.Println("Transaction for $", amt)
			if (balance - amt) < 0 {
				fmt.Println("Transaction failed!")
			} else {
				mutx.Lock()
				balance = balance - amt
				mutx.Unlock()
				fmt.Println("Transaction succeeded")
			}
			fmt.Println("Balance now $", balance)

		case status := <-tranChan:
			if status == true {

				fmt.Println("Done")
				breakPoint = true
				close(tranChan)

			}

		} // end select
	} // end for

	wg.Wait()
	fmt.Println("Final balance: $", balance)

} // end of main

func transaction(amt int) bool {
	// mutx.Lock()

	approved := false
	if (balance - amt) < 0 {
		approved = false
	} else {
		approved = true
		// mutx.Lock()
		balance = balance - amt
		// mutx.Unlock()
	}

	approvedText := "declined"
	if approved == true {
		approvedText = "approved"
	} else {

	}
	transactionNo = transactionNo + 1
	fmt.Println(transactionNo, "Transaction for $", amt, approvedText)
	fmt.Println("\tRemaining balance $", balance)

	// mutx.Unlock()
	return approved
}
