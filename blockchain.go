package main

import (
	"bytes"
	"github.com/boltdb/bolt"
	"fmt"
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
func NewBlockChain(address string) *BlockChain {
	genesisBlock := GenesisBlock(address) //创世区块
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
func GenesisBlock(address string) *Block {
	coinbase := NewCoinBaseTX(address,"我是创世区块")
	return NewBlock([]*Transaction{coinbase},[]byte{})
}

//7.添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
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
		block := NewBlock(txs,lastHash)
		//写数据  key是区块的哈希  value是区块数据转为字节流
		bucket.Put(block.Hash,block.Serialize())
		bucket.Put([]byte("LastHashKey"),block.Hash)
		//更新内存中的区块链，将tail进行更新
		bc.tail = block.Hash

		return nil
	})

}

func (bc *BlockChain)Printchain() {
	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		//假设bucket存在且有键
		b := tx.Bucket([]byte("blockBucket"))

		//从第一个key-value进行遍历，到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashkey")){
				return nil
			}

			block := Deserialize(v)
			fmt.Printf("======================区块高度：%d =======================",blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
			fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
			fmt.Printf("时间戳: %d\n", block.TimeStamp)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf( "区块数据 :%s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}

//找到指定地址的所有utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	//1.遍历区块
	//2.遍历交易
	//3.遍历output，找到和自己相关的utxo
	it := bc.NewIterator()	//创建迭代器
	for {
		block := it.Next()

		for _,tx := range block.Transactions {
			fmt.Printf("当前的交易ID %x\n",tx.TXID)
		OUTPUT:
			for i,output := range tx.TXOutputs {
				//在这里做一个过滤，将所有消耗过的outputs和当前的所即将添加的output对比一下
				//若当前的output已经消耗过，就不进行添加
				//如果当前的交易id存在于map中，说明这个交易中有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil {
					for _,j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j { //当前准备添加的output已经消耗过了
							continue OUTPUT
						}
					}
				}
					//PubKeyHash为锁定脚本
					//这个output和我们的目标地址相同，满足条件，加到返回utxo数组中
				if output.PubKeyHash == address {
					UTXO = append(UTXO,output)
				}
			}


			//如果是挖矿交易，则不用遍历直接跳过
			if !tx.IsCoinbase() {
				//4.遍历input，找到自己已经花费过的utxo集合
				//定义一个map来保存消费过的output(定义写在上面),key是这个output的交易id,value是这个交易中的索引的数组
				//map[交易id][]int64
				for _, input := range tx.TXInputs {
					//判断当前这个input和目标address是否相同，如果相同，说明这个input是该地址消耗过的output
					if input.Sig == address {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray, input.Index)
					}
				}
			} else {
				fmt.Printf("这是coinbase交易\n")
			}
		}
		if len(block.PrevHash) == 0{
			break
			fmt.Printf("区块遍历完成，退出")
		}
	}
	return UTXO
}

