package core

import (
	"fmt"
	"log"
)

type Blockchain struct {
	Blocks []*Block
}

// 创建区块链
func NewBlockchain() *Blockchain {
	genesisBlock := GenerateGenesisBlock()
	blockchain := Blockchain{}
	blockchain.AppendBlock(&genesisBlock)
	return &blockchain
}

// 新建区块
func (bc *Blockchain) SendData(data string) {
	preBlock := bc.Blocks[len(bc.Blocks) - 1]
	newBlock := GenerateNewBlock(*preBlock, data)
	bc.AppendBlock(&newBlock)
}

// 添加区块
func (bc *Blockchain) AppendBlock(newBlock *Block) {
	if len(bc.Blocks) == 0 {
		bc.Blocks = append(bc.Blocks, newBlock)
		return
	}
	if isValid(*newBlock, *bc.Blocks[len(bc.Blocks) - 1]) {
		bc.Blocks = append(bc.Blocks, newBlock)
	} else {
		log.Fatal("区块非法")
	}
}


func (bc *Blockchain) Print() {
	for _, block := range bc.Blocks {
		fmt.Println("Index:", block.Index)
		fmt.Println("TimeStamp:", block.TimeStamp)
		fmt.Println("Data:", block.Data)
		fmt.Println("Hash:", block.Hash)
		fmt.Println("PrevBlockHash:", block.PrevBlockHash)
		fmt.Println()
	}
}

// 校验区块的合法性
func isValid(newBlock, oldBlock Block) bool {
	if newBlock.Index - 1 != oldBlock.Index {
		return false
	}

	if newBlock.PrevBlockHash != oldBlock.Hash {
		return false
	}

	//if calculateHash(newBlock) != newBlock.Hash {
	//	return false
	//}
	return true
}