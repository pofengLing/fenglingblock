package main

import "fmt"

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
	return &block
}

func main()  {
	block := NewBlock("hello world",[]byte{})
	fmt.Printf("前区块哈希: %x\n",block.PrevHash)
	fmt.Printf("当前区块哈希: %x\n",block.Hash)
	fmt.Printf("区块数据: %s\n",block.Data)
}
