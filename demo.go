package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
)

type UserConfig struct {
	PrivateKey string
	NodeID     string
}

func main() {
	userConfig := UserConfig{
		"8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6",
		"047204499d849948aaffdec7ce2703f5b3",
	}

	contractBinFile, _ := ioutil.ReadFile("./demo_contracts/contract_sol_Storage.bin")
	contractBinFile = common.Hex2Bytes(string(contractBinFile))
	input := crypto.Keccak256([]byte("get()"))[0:4]
	caller := vm.AccountRef(common.HexToAddress(userConfig.NodeID))

	var levm LEVM
	levm.UserConfig = userConfig
	// 虚拟机初始化
	levm.Init()

	// 先把合约部署了
	contractAddress, _, err := levm.Deploy(caller, contractBinFile, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 再调用合约
	ret, leftgas, err := levm.Call(caller, contractAddress, input, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(ret), "left gas: ", leftgas)
}
