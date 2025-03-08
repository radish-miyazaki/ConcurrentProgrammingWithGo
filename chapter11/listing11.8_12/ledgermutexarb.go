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
	arb := NewArbitrator()
	for i := range total {
		go func(tellerID int) {
			for range 1000 {
				// 送信元と送信先をランダムに決定
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}

				// 送金処理を実行
				accounts[from].Transfer(accounts[to], 10, tellerID, arb)
			}
			fmt.Println(tellerID, "COMPLETE")
		}(i)
	}

	time.Sleep(60 * time.Second)
}

type Arbitrator struct {
	accountsInUse map[string]bool
	cond          *sync.Cond
}

func NewArbitrator() *Arbitrator {
	return &Arbitrator{
		accountsInUse: map[string]bool{},
		cond:          sync.NewCond(&sync.Mutex{}),
	}
}

func (a *Arbitrator) LockAccounts(ids ...string) {
	a.cond.L.Lock()
	defer a.cond.L.Unlock()

	for allAvailable := false; !allAvailable; {
		allAvailable = true
		for _, id := range ids {
			if a.accountsInUse[id] {
				allAvailable = false
				a.cond.Wait()
			}
		}
	}

	for _, id := range ids {
		a.accountsInUse[id] = true
	}
}

func (a *Arbitrator) UnlockAccounts(ids ...string) {
	a.cond.L.Lock()
	defer a.cond.L.Unlock()

	for _, id := range ids {
		a.accountsInUse[id] = false
	}
	a.cond.Broadcast()
}

type BankAccount struct {
	id      string
	balance int
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{id: id, balance: 100}
}

func (src *BankAccount) Transfer(
	to *BankAccount,
	amount int,
	tellerID int,
	arb *Arbitrator,
) {
	fmt.Printf("%d Locking %s and %s\n", tellerID, src.id, to.id)
	arb.LockAccounts(src.id, to.id)

	src.balance -= amount
	to.balance += amount

	arb.UnlockAccounts(src.id, to.id)

	fmt.Printf("%d Unlocked %s and %s\n", tellerID, src.id, to.id)
}
