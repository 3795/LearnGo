package core

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	// 时间戳
	Timestamp int64

	// 区块包含的数据
	Data []byte

	// 前一个区块的哈希值
	PrevBlockHash []byte

	// 区块自身的哈希值，检验区块是否有效
	Hash []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}