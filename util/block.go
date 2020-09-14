package util

import (
	"crypto/sha256"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/hashaltcoin/trx-wallet/core"
)

func GetBlockHash(block core.Block) []byte {
	rawData := block.BlockHeader.RawData

	rawDataBytes, err := proto.Marshal(rawData)
	if err != nil {
		log.Fatalln(err.Error())
	}

	h256 := sha256.New()
	h256.Write(rawDataBytes)
	blockHash := h256.Sum(nil)

	return blockHash
}
