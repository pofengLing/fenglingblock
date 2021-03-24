package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)
const walletFile = "wallet.dat"

//定义一个wallets结构，保存所有的wallet以及它的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

//创建方法
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws

}

func (ws *Wallets)CreatWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet
	ws.saveToFile()
	return address
}
//保存方法，把新建的wallet添加进去
func (ws *Wallets)saveToFile() {
	var buffer bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(walletFile, buffer.Bytes(),0600)
}

//读取文件方法，把所有的wallet读出来
func (ws *Wallets) loadFile() {
	//在读取之前要先确认文件是否存在
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		ws.WalletsMap = make(map[string]*Wallet)
		return
	}
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	//解码
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wsLocal Wallets
	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = wsLocal.WalletsMap

}

func (ws *Wallets)ListAllAddresses() []string {
	var addresses []string
	for address := range ws.WalletsMap{
		addresses = append(addresses, address)
	}
	return addresses
}