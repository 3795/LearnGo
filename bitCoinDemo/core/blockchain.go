package core

import "fmt"

type BlockChain struct {
	Blocks []*Block
}

// 添加区块到链
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// 创建初始的区块链
func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

func (bc *BlockChain) Print() {
	for _, block := range bc.Blocks {
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Prev: %x\n", string(block.PrevBlockHash))
		fmt.Println("Block Timestamp: ", block.Timestamp)
		fmt.Println()
	}
}
