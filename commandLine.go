package main

import (
	"fmt"
	"time"
)


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
		//fmt.Printf("时间戳: %x\n", block.TimeStamp)
		TimeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳: %s\n", TimeFormat)
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

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	//1.创建挖矿交易
	coinbase := NewCoinBaseTX(miner, data)
	//2.创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		fmt.Printf("无效的交易")
		return
	}
	//3.添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	//cli.bc.AddBlock([]*Transaction{tx})
	fmt.Printf("转账成功！")
}

func (cli *CLI)NewWallet() {
	//Wallet := NewWallet()
	ws := NewWallets()
	address := ws.CreatWallet()

	fmt.Printf("钱包地址: %v\n", address)

	//fmt.Printf("私钥: %v\n", Wallet.Private)
	//fmt.Printf("公钥: %v\n", Wallet.pubkey)

}
func (cli *CLI)ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddresses()
	for _, address := range addresses{
		fmt.Printf("%v\n", address)
	}
}