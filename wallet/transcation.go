package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashaltcoin/trx-wallet/common/base58"
	"github.com/hashaltcoin/trx-wallet/common/hexutil"
	"github.com/hashaltcoin/trx-wallet/service"
	"github.com/shopspring/decimal"
)

var num60 = decimal.New(1, 6)

// Send 转币
func Send(prikey, contract, to string, amount decimal.Decimal) (string, error) {
	pkBytes, _ := hex.DecodeString(prikey)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	key := (*ecdsa.PrivateKey)(privKey)

	node := service.GetGRPCClient()
	defer node.Conn.Close()
	amount6, _ := amount.Mul(num60).Float64()

	if contract == "trx" {
		return node.Transfer(key, to, int64(amount6))
	} else if contract != "" && contract != "trx" {
		data := processTransferParameter(to, int64(amount6))
		return node.TransferContract(key, contract, data)
	}

	return "", fmt.Errorf("the type %s not support now", contract)
}

// 处理合约转账参数
func processTransferParameter(to string, amount int64) (data []byte) {
	methodID, _ := hexutil.Decode("a9059cbb")
	addr, _ := base58.DecodeCheck(to)
	paddedAddress := common.LeftPadBytes(addr, 32)
	amountBig := new(big.Int).SetInt64(amount)
	paddedAmount := common.LeftPadBytes(amountBig.Bytes(), 32)
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return
}

// 处理合约获取余额
func processBalanceOfData(trc20 []byte) (amount int64) {
	if len(trc20) >= 32 {
		amount = new(big.Int).SetBytes(common.TrimLeftZeroes(trc20[0:32])).Int64()
	}
	return
}

// 处理合约获取余额参数
func processBalanceOfParameter(addr string) (data []byte) {
	methodID, _ := hexutil.Decode("70a08231")
	add, _ := base58.DecodeCheck(addr)
	paddedAddress := common.LeftPadBytes(add, 32)
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	return
}
