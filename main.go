package main

func main() {
	//产生区块链
	bc := NewBlockChain("魄风灵")
	cli := CLI{bc}
	cli.Run()
	//bc.AddBlock("我是第002个区块")
	//bc.AddBlock("我是第003个区块")
	//
	//
	////调用迭代器返回每一个区块数据
	//it := bc.NewIterator()
	//for {
	//	//返回区块，左移
	//	block := it.Next()
	//
	//	fmt.Printf("==================================\n")
	//	fmt.Printf("前区块哈希: %x\n", block.PrevHash)
	//	fmt.Printf("当前区块哈希: %x\n", block.Hash)
	//	fmt.Printf("区块数据: %s\n", block.Data)
	//	if len(block.PrevHash) == 0 {
	//		fmt.Println("\n区块链遍历结束")
	//		break
	//	}
	//}

}
