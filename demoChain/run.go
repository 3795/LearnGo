package main

import "demoChain/core"

func main() {
	bc := core.NewBlockchain()
	bc.SendData("send 1 BTC to Jacky")
	bc.SendData("send 1 EOS to Jacky")
	bc.Print()
}
