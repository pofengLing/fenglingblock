package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 50

//1.定义交易结构
type Transaction struct {
	TXID []byte				//交易ID
	TXInputs []TXInput		//交易输入数组
	TXOutputs []TXOutput	//交易输出数组
}

//定义交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本，用地址模拟
	Sig string
}

//定义交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，用地址模拟
	PubKeyHash string

}

//设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]

}

//判断当前的交易是否是挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	//1.交易input只有一个
	//2.交易ID为空
	//3.交易的index为-1
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
			return true
		}
	return false

}


//2.提供创建交易方法
//先实现最简单的挖矿交易(coinbase)
//挖矿交易的特点：1.只有一个input 2.无需引用交易id  3.无需引用index
func NewCoinBaseTX(address string, data string) *Transaction{
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自行填写，一般填矿池的名字
	input := TXInput{[]byte{},-1,data}
	output := TXOutput{reward,address}
	//对于挖矿交易，只有一个input和一个output
	tx := Transaction{[]byte{},[]TXInput{input},[]TXOutput{output}}
	//设置交易ID
	tx.SetHash()
	return &tx
}

//创建普通的转账交易
//1.找到最合适的UTXO集合 map[string][]uint64
//2.将这些UTXO逐一转成input
//3.创建outputs
//4.如果有零钱，就进行找零
func NewTransaction(from, to string,amount float64,bc *BlockChain) *Transaction {
	//找到最合理的UTXO集合 map[string][]uint64
	utxos, resValue := bc.FindNeedUTXOs(from, amount)

	if resValue < amount {
		fmt.Printf("余额不足，交易失败")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//将这些UTXO逐一转成input
	for id, indexArray := range utxos{
		for _, i := range indexArray {
			input := TXInput{[]byte(id),int64(i),from}
			inputs = append(inputs, input)
		}
	}

	//创建outputs
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	if resValue > amount {
		//找零
		outputs = append(outputs, TXOutput{resValue-amount, from})
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}


//4.根据交易调整程序