package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	//这里的pubkey不存储原生的公钥，而是直接存储punkey结构体中的X、Y拼接而成的字符串，在校验端重新拆分
	pubkey []byte
}

//创建钱包
func NewWallet() *Wallet{
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}
	//生成公钥
	pubKeyOrig := privateKey.PublicKey
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)
	return &Wallet{privateKey, pubKey}
}


//生成地址
func (w *Wallet) NewAddress() string {
	pubkey := w.pubkey
	hash := sha256.Sum256(pubkey)
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	//返回rip160的哈希结果
	rip160HashValue := rip160hasher.Sum(nil)
	//版本号 1字节
	version := byte(00)
	payload := append([]byte{version}, rip160HashValue...)

	//checksum
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	//前4字节校验码
	checkCode := hash2[:4]
	//25字节数据
	payload = append(payload, checkCode...)

	address := base64.StdEncoding.EncodeToString(payload)
	return address

}