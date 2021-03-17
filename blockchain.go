package main

import (
	"github.com/boltdb/bolt"
	"log"
)

//4.引入区块链
type BlockChain struct {
	//定义一个区块链数组
	//改写区块链，用数据库取代[]*Block
	//blocks []*Block;
	db *bolt.DB
	//由于存储在数据库中，需将最后一个块的hash单独保存下来
	tail []byte //记录最后一个块的哈希值
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

//5.定义一个区块链
func NewBlockChain() *BlockChain{
	genesisBlock := GenesisBlock()//创世区块
	//原始方法
	//return &BlockChain{
	//	blocks: []*Block{genesisBlock}, //将创世区块添加到区块链中
	//}

	var lastHash []byte //最后一个区块的哈希，从数据库中读取出来的

	//数据库方法
	db,err := bolt.Open(blockChainDb,0600,nil)
	//defer db.Close()
	if err != nil{
		log.Panic("打开数据库失败")
	}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有bucket，创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket(blockBucket)失败")
			}



			//写数据  key是区块的哈希  value是区块数据转为字节流
			bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"),genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		}else {
			//有bucket，则先读
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db,lastHash}
}

//6.生成创世区块
func GenesisBlock() *Block{
	return NewBlock("我是创世区块",[]byte{})
}

//7.添加区块
func (bc *BlockChain) AddBlock(data string) {
	//原始方法
	//取出前一个区块
	//lastBlock := bc.blocks[len(bc.blocks)-1]
	//计算前一个区块的的哈希值
	//prevHash := lastBlock.Hash
	//a.创建新的区块
	//block := NewBlock(data,prevHash)
	//b.将新的区块添加到区块链中
	//bc.blocks = append(bc.blocks,block)

	//数据库方法
	//先拿到数据库
	db := bc.db //区块链数据库
	lastHash := bc.tail //最后一个区块的哈希


	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket不应该为空，请检查")
		}
		block := NewBlock(data,lastHash)
		//写数据  key是区块的哈希  value是区块数据转为字节流
		bucket.Put(block.Hash,block.Serialize())
		bucket.Put([]byte("LastHashKey"),block.Hash)
		//更新内存中的区块链，将tail进行更新
		bc.tail = block.Hash

		return nil
	})

}
