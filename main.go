package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/samuelgoes/Web3GoTokio/contracts"
	"log"
	"math/big"
	"time"
)

func main() {
	// Use Ganache node:
	/*
		cl, err := ethclient.Dial("http://127.0.0.1:7545")
		if err != nil {
			log.Fatalf("error dialing eth client: %v", err)
		}
	*/

	// Use Infura:
	infura := "https://goerli.infura.io/v3/76273ae5f1af44cabc486caeb2fa28aa"
	cl, err := ethclient.Dial(infura)

	defer cl.Close()

	hexPrivKey := "f713f261307e511dcb0030a6d5c7b4022ae1c8f7deb614ea2f2a1f6b9d8ed738"
	key, err := crypto.HexToECDSA(hexPrivKey)
	if err != nil {
		log.Fatalf("Private key is not OK. %v", err)
	}

	addr := common.HexToAddress("0x227d0f9A88Dd26Fa05e83Ac2008082Ff53A2541d")
	ctx := context.Background()

	// Retrieve a block by number
	block, err := cl.BlockByNumber(ctx, big.NewInt(37))
	if err != nil {
		log.Printf("error getting block number: %v", err)
	} else {
		log.Printf("Block: Transactaions: %v", block.Transactions())
	}

	// Get Balance of an account (nil means at newest block)
	balance, err := cl.BalanceAt(ctx, addr, nil)
	if err != nil {
		log.Fatalf("error getting balance: %v", err)
	}
	log.Printf("Balance: %v", balance)

	// Get sync progress of the node. If nil, the node is not syncing
	progress, err := cl.SyncProgress(ctx)
	if err != nil {
		log.Fatalf("error getting balance: %v", err)
	}

	if progress != nil {
		log.Printf("Progress: %v", progress)
	}

	// ****************** DEPLOY ******************

	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := cl.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("error reading next nonce")
	}

	gasPrice, err := cl.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("error reading suggest gas price")
	}

	chainID, err := cl.ChainID(context.Background())
	if err != nil {
		log.Fatalf("unable to get chainID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Fatalf("unable to build new transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	// Deploy Smart Contract
	address, txdc, con, err := contracts.DeployContracts(auth, cl)
	if err != nil {
		log.Fatalf("unable to deploy smart contract. %v", err)
	}

	waitForBlock(cl, txdc)
	log.Printf("Smart Contract desplegado satisfactoriamente. Address: %s", address.Hex())

	/*
		// Load Smart Contract
		contractAddresss := common.HexToAddress("0x289411fd5C6E8f8e27321EA30Cb988C4c5585509")
		con, err := contracts.NewContracts(contractAddresss, cl)
		if err != nil {
			log.Fatalf("unable to load smart contract")
		}
		log.Printf("Smart Contract cargado satisfactoriamente.")
	*/

	nonce, err = cl.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("error reading next nonce")
	}
	auth.Nonce = big.NewInt(int64(nonce))

	addr1 := common.HexToAddress("0xAa29b832234876114daf40f07398F1Bc37d3963c")
	amount := big.NewInt(1800000000000)

	tx1, err := con.Transfer(auth, addr1, amount)
	if err != nil {
		log.Fatalf("unable to call store message function. %v", err)
	}
	waitForBlock(cl, tx1)
	log.Printf("Transferencia realizada satisfactoriamente")

	gas := tx1.Gas()
	log.Printf("Gas usado en la Tx: %v", gas)

	symbol, err := con.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("unable to call Symbol function. Err: %v", err)
	}
	log.Printf("Este es el símbolo del Smart Contract TokioToken: %s", symbol)

	balanceOf, err := con.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		log.Fatalf("unable to call BalanceOf function. Err: %v", err)
	}
	log.Printf("Este es el balance de la cuenta principal - ERC-20: %v", balanceOf)

	balanceOf, err = con.BalanceOf(&bind.CallOpts{}, addr1)
	if err != nil {
		log.Fatalf("unable to call BalanceOf function. Err: %v", err)
	}
	log.Printf("Este es el balance de la cuenta Tokio - ERC-20: %v", balanceOf)
}

func waitForBlock(cl *ethclient.Client, tx *types.Transaction) {
	for true {
		_, err := cl.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Printf("Esperando 5seg")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}
