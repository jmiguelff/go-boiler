package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	conn, err := ethclient.Dial("https://mainnet.infura.io/v3/e65cecf22d7f4e29bbc0a4dc424a7c2c")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := NewTobyTheGoat(common.HexToAddress("0x2484A62c2D3C6980Dc57D9aa02305C2F7523Dc0d"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	fd, err := os.OpenFile("tokens.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer fd.Close()

	for i := 0; i < 9999; i++ {
		id := big.NewInt(int64(i))
		tok, err := contract.TokenURI(&bind.CallOpts{}, id)
		if err != nil {
			fmt.Printf("Failed to get token URI (%d): %v", i, err)
		} else {
			fd.WriteString(tok)
			fd.WriteString("\n")
			fmt.Println(tok)
		}
	}

	//sup, _ := contract.TotalSupply(&bind.CallOpts{})
	//amt, err := contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0x5F04C0D06CEB108Bf2BCc8C3fBdfDD4b061085a8"))
	//if err != nil {
	//	log.Fatalf("Failed to get total supply: %v", err)
	//}
	//fmt.Println(sup)
}
