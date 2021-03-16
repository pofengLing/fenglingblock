package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

//0.定义结构
type Block struct {
	//1.版本号
	Version uint64
	//2.前区块哈希
	PrevHash []byte
	//3.Merkel根
	MerkelRoot []byte
	//4.时间戳
	TimeStamp uint64
	//5.难度值
	Difficulty uint64
	//6.随机数（挖矿要找的数据）
	Nonce uint64

	//a.当前区块哈希
	Hash []byte
	//b.数据
	Data []byte
}

//辅助函数，将uint转为[]byte
func Uint64ToByte(num uint64) []byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}

	return buffer.Bytes()

	//2.
	//var Bytes []byte
	//Bytes = make([]byte,10)
	//binary.BigEndian.PutUint64(Bytes,num)
	//return Bytes
}

//2.创建区块
func NewBlock(data string,prevBlockHash []byte) *Block{
	block := Block{
		Version: 00,
		PrevHash: prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		Difficulty: 0,  //随便填写的无效值
		Nonce: 0,
		Hash: []byte{},//先填空，之后计算
		Data: []byte(data),
	}
	//block.SetHash() 此种方法是固定计算，应该改为使用pow的run函数动态计算
	pow := NewProofOfWork(&block)
	//查找随机数不停进行哈希运算
	hash,nonce := pow.Run()

	//根据运算结果对区块数据进行更新
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//3.生成哈希
func (block *Block) SetHash() []byte{
	var blockInfo []byte
	//1.拼装数据
	/*
	blockInfo = append(blockInfo,Uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo,block.PrevHash...)
	blockInfo = append(blockInfo,block.MerkelRoot...)
	blockInfo = append(blockInfo,Uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo,Uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo,Uint64ToByte(block.Nonce)...)
	blockInfo = append(blockInfo,block.Data...)
	 */ //这种方法比较直接

	 //join优化方法
	 tmp := [][]byte{
	 	Uint64ToByte(block.Version),
	 	block.PrevHash,
	 	block.MerkelRoot,
	 	Uint64ToByte(block.TimeStamp),
	 	Uint64ToByte(block.Difficulty),
	 	Uint64ToByte(block.Nonce),
	 	block.Data,
	 }
	 blockInfo = bytes.Join(tmp,[]byte{})

	//2.sha256
	//func Sum256(data []byte) [size]byte{}
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
	return block.Hash
}
