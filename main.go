package main

import (
	"crypto/sha256"
	"fmt"
)

//0.定义结构
type Block struct {
	//1.前区块哈希
	PrevHash []byte
	//2.当前区块哈希
	Hash []byte
	//3.数据
	Data []byte
}

//2.创建区块
func NewBlock(data string,prevBlockHash []byte) *Block{
	block := Block{
		PrevHash:prevBlockHash,
		Hash: []byte{},//先填空，之后计算
		Data: []byte(data),
	}
	block.SetHash()
	return &block
}

//3.生成哈希
func (block *Block) SetHash() []byte{
	//1.拼装数据
	blockInfo := append(block.PrevHash,block.Data...)
	//2.sha256
	//func Sum256(data []byte) [size]byte{}
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
	return block.Hash
}

//4.引入区块链
type BlockChain struct {
	//定义一个区块链数组
	blocks []*Block;
}

//5.定义一个区块链
func NewBlockChain() *BlockChain{
	genesisBlock := GenesisBlock()//创世区块
	return &BlockChain{
		blocks: []*Block{genesisBlock}, //将创世区块添加到区块链中
	}
}

//6.生成创世区块
func GenesisBlock() *Block{
	return NewBlock("我是创世区块",[]byte{})
}


func main()  {
	//产生区块链
	bc := NewBlockChain()
	for i,block := range bc.blocks {
		fmt.Printf("当前区块高度: %d\n",i)
		fmt.Printf("前区块哈希: %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希: %x\n", block.Hash)
		fmt.Printf("区块数据: %s\n", block.Data)
	}
}
