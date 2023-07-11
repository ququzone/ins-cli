package ins

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"

	"github.com/ququzone/ins-cli/pkg/contracts"
)

var (
	duration = big.NewInt(365 * 24 * 60 * 60)
	secret   = [32]byte{}

	errCommit   = errors.New("commit commitment error")
	errRegister = errors.New("register name error")
)

func init() {
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte("secret"))
	hash := sha.Sum(nil)
	copy(secret[:], hash)
}

func nameHash(name string) (hash [32]byte, err error) {
	var ioNode []byte
	ioNode, err = hex.DecodeString("b2b692c69df4aa3b0a24634d20a3ba1b44c3299d09d6c4377577e20b09e68395")
	if err != nil {
		return
	}

	sha := sha3.NewLegacyKeccak256()
	if _, err = sha.Write(ioNode); err != nil {
		return
	}
	nameSha := sha3.NewLegacyKeccak256()
	if _, err = nameSha.Write([]byte(name)); err != nil {
		return
	}
	nameHash := nameSha.Sum(nil)
	if _, err = sha.Write(nameHash); err != nil {
		return
	}
	sha.Sum(hash[:0])
	return
}

func Register(rpc, controllerAddr, resolverAddr, privateKey, owner, name string) error {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}

	if strings.EqualFold("0x", privateKey[:2]) {
		privateKey = privateKey[2:]
	}
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}
	ownerAddr := crypto.PubkeyToAddress(key.PublicKey)
	if owner != "" {
		ownerAddr = common.HexToAddress(owner)
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		return err
	}
	controller, err := contracts.NewIOTXRegistrarController(common.HexToAddress(controllerAddr), client)
	if err != nil {
		return err
	}

	hash, err := nameHash(name)
	if err != nil {
		return err
	}
	resolverABI, err := abi.JSON(strings.NewReader(`[
		{
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "a",
          "type": "address"
        }
      ],
      "name": "setAddr",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
	]`))
	if err != nil {
		return err
	}
	data, err := resolverABI.Pack("setAddr", hash, ownerAddr)
	if err != nil {
		return err
	}

	commitment, err := controller.MakeCommitment(
		nil,
		name,
		ownerAddr,
		duration,
		secret,
		common.HexToAddress(resolverAddr),
		[][]byte{data},
		true,
		0,
	)
	if err != nil {
		return err
	}
	fmt.Printf("commit commitment: %s ...\n", hex.EncodeToString(commitment[:]))

	tx, err := controller.Commit(transactor, commitment)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return err
	}
	if receipt.Status != 1 {
		return errCommit
	}

	fmt.Printf("sleep for activation commitment ...\n")
	time.Sleep(60 * time.Second)
	price, err := controller.RentPrice(nil, name, duration)
	if err != nil {
		return err
	}
	transactor.Value = new(big.Int).Add(price.Base, price.Premium)
	tx, err = controller.Register(transactor,
		name,
		ownerAddr,
		duration,
		secret,
		common.HexToAddress(resolverAddr),
		[][]byte{data},
		true,
		0)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return err
	}
	if receipt.Status != 1 {
		return errRegister
	}
	fmt.Printf("Register %s txHash: %s\n", name, tx.Hash().String())

	return nil
}
