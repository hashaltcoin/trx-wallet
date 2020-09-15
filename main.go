package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/btcsuite/btcd/btcec"
	"github.com/hashaltcoin/trx-wallet/common/base58"
	"github.com/hashaltcoin/trx-wallet/common/crypto"
	"github.com/hashaltcoin/trx-wallet/config"
	pb "github.com/hashaltcoin/trx-wallet/tron"
	"github.com/hashaltcoin/trx-wallet/wallet"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

var log *config.Logger

type TronWallet struct {
}

func (tw *TronWallet) GetBalance(ctx context.Context, token *pb.Token) (*pb.TokenBalance, error) {
	pkBytes, _ := hex.DecodeString(token.Address)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	re := (*ecdsa.PrivateKey)(privKey)
	addr := base58.EncodeCheck(crypto.PubkeyToAddress(re.PublicKey).Bytes())

	balance := int64(0)
	if token.Code == "trx" {
		balance = wallet.GetBalance(addr)
	} else if token.Code != "trx" && token.Code != "" {
		balance = wallet.GetTokenBalance(re, addr, token.Code)
	}

	return &pb.TokenBalance{Amount: balance}, nil
}

func (tw *TronWallet) Transfer(ctx context.Context, in *pb.SendInfo) (*pb.TransferResult, error) {
	log.Info("Received: %d", in.Index)
	txid, err := wallet.Send(in.Private, "trx", in.To, decimal.New(in.Amount, -6))

	return &pb.TransferResult{Code: 2000, Txid: txid}, err
}

func (tw *TronWallet) TransferToken(ctx context.Context, in *pb.SendTokenInfo) (*pb.TransferResult, error) {
	log.Info("Received: %x", in)
	txid, err := wallet.Send(in.Private, in.Contract, in.To, decimal.New(in.Amount, -6))

	return &pb.TransferResult{Code: 2000, Txid: txid}, err
}

func (tw *TronWallet) GetWallet(ctx context.Context, in *pb.WalletIndex) (*pb.WalletInfo, error) {
	log.Info("Received: %d", in.Index)
	pri, addr, err := wallet.CreateAccount()
	return &pb.WalletInfo{Private: pri, Address: addr}, err
}

func main() {
	log = config.NewLogger("INFO", "Main", "TRX-Wallet")
	cfg, _ := config.GetConfig()

	log.Info("Server: %s:%d", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Error("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterContractServer(s, &TronWallet{})
	if err := s.Serve(lis); err != nil {
		log.Error("failed to serve: %v", err)
	}
}
