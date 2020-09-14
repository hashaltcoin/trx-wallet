package main

import (
	"context"
	"fmt"
	"net"

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
	balance := int64(0)
	if token.Code == "trx" {
		balance = wallet.GetBalance(token.Address)
	} else if token.Code != "trx" && token.Code != "" {
		balance = wallet.GetTokenBalance(token.Address, token.Code)
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
