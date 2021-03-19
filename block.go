package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
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

	//a.当前区块哈希  正常比特币区块中没有当前区块哈希
	Hash []byte
	//b.数据
	//Data []byte
	//真实交易数组
	Transactions []*Transaction
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

//辅助函数，将数据转换为字节流
//序列化函数
func (block *Block) Serialize() []byte{
	//定义存储编码的数据的buffer
	var buffer bytes.Buffer
	//使用god进行序列化（编码）得到字节流
	//1.定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	//2.使用编码器进行编码
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("编码出错")
	}

	//输出编码后的数据
	//fmt.Printf("编码后的数据： %v\n",buffer.Bytes())

	return buffer.Bytes()
}


//反序列化函数
func Deserialize(data []byte) Block {
	//定义解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block
	//使用解码器进行解码
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错")
	}

	return block
}

//2.创建区块
func NewBlock(txs []*Transaction,prevBlockHash []byte) *Block {
	block := Block{
		Version: 		00,
		PrevHash: 		prevBlockHash,
		MerkelRoot: 	[]byte{},
		TimeStamp: 		uint64(time.Now().Unix()),
		Difficulty: 	0,  //随便填写的无效值
		Nonce: 			0,
		Hash: 			[]byte{},//先填空，之后计算
		Transactions: 	txs,
	}
	block.MerkelRoot = block.MakeMerKelRoot()
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
	 	 //block.Data,
	 }
	 blockInfo = bytes.Join(tmp,[]byte{})

	//2.sha256
	//func Sum256(data []byte) [size]byte{}
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
	return block.Hash
}

//生成merkel根  模拟实现，只做简单拼接 实际上应该用二叉树实现
func (block *Block) MakeMerKelRoot() []byte{

	return []byte{}
}
