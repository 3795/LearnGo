package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// 区块机构
type Block struct {
	// 区块编号
	Index int64

	// 区块时间戳
	TimeStamp int64

	// 上一个区块哈希值
	PrevBlockHash string

	// 当前区块哈希值
	Hash string

	// 区块数据
	Data string
}

// 计算区块的Hash值
func calculateHash(block Block) string {
	blockData := string(block.Index) + string(block.TimeStamp) +
		string(block.PrevBlockHash) + block.Data
	hashInBytes :=sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hashInBytes[:])
}

// 生成新的区块
func GenerateNewBlock(preBlock Block, data string) Block {
	newBlock := Block{}
	newBlock.Index = preBlock.Index + 1
	newBlock.PrevBlockHash = preBlock.Hash
	newBlock.TimeStamp = time.Now().Unix()
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Data = data
	return newBlock
}

// 生成创世区块
func GenerateGenesisBlock() Block {
	preBlock := Block{}
	preBlock.Index = -1
	preBlock.Hash = ""
	return GenerateNewBlock(preBlock, "Genesis Block")
}