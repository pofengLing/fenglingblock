package main

import "fmt"


func (cli *CLI) AddBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Printf("添加区块成功！！！\n")
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	//调用迭代器返回每一个区块数据
	it := bc.NewIterator()
	for {
		//返回区块，左移
		block := it.Next()

		fmt.Printf("==================================\n")
		fmt.Printf("版本号: %x\n", block.Version)
		fmt.Printf("前区块哈希: %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希: %x\n", block.Hash)
		fmt.Printf("merkle根: %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %x\n", block.TimeStamp)
		fmt.Printf("难度值: %x\n", block.Difficulty)
		fmt.Printf("随机数: %x\n", block.Nonce)
		fmt.Printf("区块数据: %s\n", block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Println("\n区块链遍历结束")
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	utxos := cli.bc.FindUTXOs(address)
	total := 0.0
	for _,utxo := range utxos{
		total += utxo.Value
	}
	fmt.Printf("\"%s\"余额为：%f\n",address,total)
}
