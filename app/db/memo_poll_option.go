package db

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type MemoPollOption struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte
	PollTxHash []byte `gorm:"index:poll_tx_hash"`
	Option     string `gorm:"size:500"`
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoPollOption) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo poll option", result.Error)
	}
	return nil
}

func (m MemoPollOption) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo poll option", err).Print()
		return ""
	}
	return hash.String()
}

func GetMemoPollOption(txHash []byte) (*MemoPollOption, error) {
	var memoPollOption MemoPollOption
	err := find(&memoPollOption, MemoPollOption{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo poll option", err)
	}
	return &memoPollOption, nil
}
