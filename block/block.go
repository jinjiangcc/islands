package block

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jiangjincc/islands/utils"

	"github.com/jedib0t/go-pretty/table"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Txs           []*Transaction
	Hash          []byte
	Height        int64
	Nonce         int64
}

// 创建新的区块
func NewBlock(txs []*Transaction, height int64, PrevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().UTC().Unix(),
		PrevBlockHash: PrevBlockHash,
		Txs:           txs,
		Height:        height,
	}

	// 工作量证明
	powIns := NewProofOfWork(block)
	hash, nonce := powIns.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func CreateGenesisBlock(data []*Transaction) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

// TODO 使用protobuff
func (b *Block) Serialize() []byte {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (b *Block) HashTransaction() []byte {
	//var (
	//	txs  [][]byte
	//	hash [32]byte
	//)
	//
	//// 只对hash进行计算
	//for i := 0; i < len(b.Txs); i++ {
	//	txs = append(txs, b.Txs[i].TxHash)
	//}
	//hash = sha256.Sum256(bytes.Join(txs, []byte{}))
	var transaction [][]byte
	for _, tx := range b.Txs {
		transaction = append(transaction, tx.Serialize())
	}
	mTree := NewMerkleTree(transaction)
	return mTree.RootNode.Data
}

func (b *Block) PrintBlock() {

	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	tb.AppendHeader(table.Row{"内容", "区块信息"})

	for _, v := range b.Txs {
		fmt.Println("交易信息:")

		for _, vv := range v.In {
			fmt.Println("交易输出", hex.EncodeToString(utils.Ripemd160Hash(vv.PublicKey)), vv.Vout, hex.EncodeToString(vv.TxHash))
		}
		for _, vv := range v.Out {
			fmt.Println("未花费", vv.Value, hex.EncodeToString(vv.PublicKey))
		}
	}

	tb.AppendRows([]table.Row{
		{"Height", b.Height},
		{"Txs", hex.EncodeToString(b.HashTransaction())},
		{"Timestamp", time.Unix(b.Timestamp, 0).Format("2006-01-02 15:04:05")},
		{"Nonce", b.Nonce},
		{"Hash", byteForString(b.Hash)},
		{"PrevHash", byteForString(b.PrevBlockHash)},
	})
	tb.SetStyle(table.StyleDefault)
	tb.Render()
}

func UnSerialize(bt []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(bt))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func byteForString(b []byte) string {
	return fmt.Sprintf("%x", b)
}
