package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	core "github.com/stevewooo/LonelyEVM/levm"
)

func main() {
	// 获取编译后的合约数据。hex格式
	contractBinFile, _ := ioutil.ReadFile("./demo_contracts/contract_sol_Hello.bin")

	YOUR_NODE_ID := "YOUE_NODE_ID"

	// 初始化
	var levm core.LEVM
	levm.Init()

	// 创建虚拟机
	levm.CreateEVM()

	// 1. 部署合约到state上
	contractAddress, _, err := levm.Deploy(YOUR_NODE_ID, string(contractBinFile), levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取调用的函数地址。函数名和参数名组成的字符串，哈希后的8个位
	input := crypto.Keccak256([]byte("get()"))[0:4]
	fmt.Println("Calling Address:", contractAddress)
	fmt.Println("Calling Point:", common.Bytes2Hex(input))

	// 2. 调用合约
	ret, _, err := levm.Call(YOUR_NODE_ID, contractAddress, common.Bytes2Hex(input), levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ret)
}
