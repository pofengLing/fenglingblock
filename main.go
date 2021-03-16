package main

import (
	"fmt"
)

func main() {
	//产生区块链
	bc := NewBlockChain()
	bc.AddBlock("我是第002个区块")
	bc.AddBlock("我是第003个区块")
	for i,block := range bc.blocks {
		fmt.Printf("当前区块高度: %d\n",i)
		fmt.Printf("前区块哈希: %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希: %x\n", block.Hash)
		fmt.Printf("区块数据: %s\n", block.Data)
	}
}
