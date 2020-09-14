package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/hashaltcoin/trx-wallet/common/base58"
	"github.com/hashaltcoin/trx-wallet/common/crypto"
	"github.com/hashaltcoin/trx-wallet/service"
)

func CreateAccount() (string, string, error) {
	re, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	addr := base58.EncodeCheck(crypto.PubkeyToAddress(re.PublicKey).Bytes())
	prikey := crypto.PrikeyToHexString(re)

	return prikey, addr, nil
}

func GetTokenBalance(pri, contract string) int64 {
	pkBytes, _ := hex.DecodeString(pri)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	re := (*ecdsa.PrivateKey)(privKey)
	addr := base58.EncodeCheck(crypto.PubkeyToAddress(re.PublicKey).Bytes())

	client := service.GetGRPCClient()
	if rsp, err := client.GetConstantResultOfContract(re, contract, processBalanceOfParameter(addr)); err != nil {
		return 0
	} else {
		return processBalanceOfData(rsp[0])
	}
}

func GetBalance(address string) int64 {
	client := service.GetGRPCClient()
	if acc, err := client.GetAccount(address); err == nil {
		return acc.Balance
	}

	return 0
}
