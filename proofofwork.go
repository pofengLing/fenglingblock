package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//1.定义一个工作量证明POW

type ProofOfWork struct {
	//a.block
	block *Block
	//b.目标值
	//一个非常大的数，有很多方法
	target *big.Int
}

//2.提供创建POW的函数
//NewProofOfWork(参数)
func NewProofOfWork(block *Block) *ProofOfWork{
	pow := ProofOfWork{
		block : block,
	}

	//难度值此时是string类型，需转换
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	//引入辅助变量，目的是将上面的难度值转换为big.Int
	tmpInt := big.Int{}
	//将难度值赋值给tmpInt,指定16进制
	tmpInt.SetString(targetStr,16)
	pow.target = &tmpInt
	return &pow
}


//3.提供不断计算哈希的函数
//Run()
func (pow *ProofOfWork) Run() ([]byte,uint64){ //[]byte是hash，uint64是nonce

	var nonce uint64
	block := pow.block
	var hash [32]byte
	//1.拼装数据（区块的数据+nonce）

	for {
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			block.Data,
		}
		blockInfo := bytes.Join(tmp,[]byte{})

		//2.做哈希运算
		hash = sha256.Sum256(blockInfo)

		//3.与pow中的目标值target进行比较
		tmpInt := big.Int{}
		//将计算出来的hash转换为bigInt
		tmpInt.SetBytes(hash[:])
		//进行比较当前哈希值和目标哈希值
			// Cmp compares x and y and returns:
			//
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
		if tmpInt.Cmp(pow.target) == -1 {
			//a.找到了，返回
			fmt.Printf("挖矿成功！hash: %x\tnonce: %x\n", hash, nonce)
			//break
			return hash[:], nonce

		}else{
			//b.未找到，nonce+1继续找
			nonce++
		}
	}

}
//4.提供校验函数
//
//Isvalid()