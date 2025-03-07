package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	accounts := []*BankAccount{
		NewBankAccount("Sam"),
		NewBankAccount("Paul"),
		NewBankAccount("Amy"),
		NewBankAccount("Mia"),
	}

	total := len(accounts)
	for i := range total {
		go func(eID int) {
			for range 1000 {
				// 送信元と送信先をランダムに決定
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}

				// 送金処理を実行
				accounts[from].Transfer(accounts[to], 10, eID)
			}
			fmt.Println(eID, "COMPLETE")
		}(i)
	}

	time.Sleep(40 * time.Second)
}

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id: id, balance: 100, mutex: sync.Mutex{},
	}
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, exID int) {
	fmt.Printf("%d Locking %s's account\n", exID, src.id)
	src.mutex.Lock()

	fmt.Printf("%d Locking %s's account\n", exID, to.id)
	to.mutex.Lock()

	src.balance -= amount
	to.balance += amount
	to.mutex.Unlock()
	src.mutex.Unlock()

	fmt.Printf("%d Unlocked %s and %s\n", exID, src.id, to.id)
}
