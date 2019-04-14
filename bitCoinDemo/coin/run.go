package main

import (
	"bitCoinDemo/core"
)

func main() {
	bc := core.NewBlockChain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	bc.Print()
}
