package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward = 12.5

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
	if len(tx.TXInputs) == 1 {
		input := tx.TXInputs[0]
		if !bytes.Equal(input.TXid,[]byte{}) && input.Index != -1{
			return false
		}
	}
	return true
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
//3.创建挖矿交易（没有输入）
//4.根据交易调整程序