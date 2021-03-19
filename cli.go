package main

import (
	"fmt"
	"os"
)

//用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA	"添加区块"
	printChain              "正向打印区块链"
	getBalance --address ADDRESS  "获取指定地址的余额"
`

//接收参数的动作放到一个函数中
func(cli *CLI) Run() {
	//1.得到所有命令
	//block printChain
	//block addBlock --data "helloWorld"
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}

	//2.分析命令
	cmd := args[1]
	switch cmd {
		case "addBlock":
			//添加区块
			fmt.Printf("添加区块\n")
			//确保数据有效
			if len(args) == 4 && args[2] == "--data"{
				//a.获取数据
				data := args[3]
				//b.使用bc添加区块addBlock
				cli.AddBlock(data)
			} else {
				fmt.Printf("添加区块命令参数不正确，请检查")
			}

		case "printChain":
			//打印区块
			fmt.Printf("打印区块\n")
			cli.PrintBlockChain()
		case "getBalance":
			//获取余额
			fmt.Printf("获取余额\n")
			//确保数据有效
			if len(args) == 4 && args[2] == "--address" {
				address := args[3]
				cli.getBalance(address)
			}
		default:
			fmt.Printf("无效的命令，请检查\n")
			fmt.Printf(Usage)
	}

}

