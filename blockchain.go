package main

//4.引入区块链
type BlockChain struct {
	//定义一个区块链数组
	blocks []*Block;
}

//5.定义一个区块链
func NewBlockChain() *BlockChain{
	genesisBlock := GenesisBlock()//创世区块
	return &BlockChain{
		blocks: []*Block{genesisBlock}, //将创世区块添加到区块链中
	}
}

//6.生成创世区块
func GenesisBlock() *Block{
	return NewBlock("我是创世区块",[]byte{})
}

//7.添加区块
func (bc *BlockChain)AddBlock(data string,){
	//取出前一个区块
	lastBlock := bc.blocks[len(bc.blocks)-1]
	//计算前一个区块的的哈希值
	prevHash := lastBlock.Hash
	//a.创建新的区块
	block := NewBlock(data,prevHash)
	//b.将新的区块添加到区块链中
	bc.blocks = append(bc.blocks,block)
}
