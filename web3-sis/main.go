package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	conn, err := ethclient.Dial("https://rinkeby.infura.io/v3/e65cecf22d7f4e29bbc0a4dc424a7c2c")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := NewSis(common.HexToAddress("0xdDe41bDa34481e9c73A06531C9f4E54DDAc2d4a2"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	sale, err := contract.SaleIsOpen(&bind.CallOpts{})
	if err != nil {
		fmt.Printf("Failed to get sale state: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(sale)
	os.Exit(0)

	//sup, _ := contract.TotalSupply(&bind.CallOpts{})
	//amt, err := contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0x5F04C0D06CEB108Bf2BCc8C3fBdfDD4b061085a8"))
	//if err != nil {
	//	log.Fatalf("Failed to get total supply: %v", err)
	//}
	//fmt.Println(sup)
}
