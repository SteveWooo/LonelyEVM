package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/core/vm"
	core "github.com/stevewooo/LonelyEVM/levm"
)

func main() {
	// 获取编译后的合约数据。hex格式
	contractBinFile, _ := ioutil.ReadFile("./demo_contracts/contract_sol_Hello.bin")
	contractBinFile = common.Hex2Bytes(string(contractBinFile))

	// 调用发起者
	caller := vm.AccountRef(common.HexToAddress("YOUE_NODE_ID"))

	// 虚拟机初始化
	var levm core.LEVM
	levm.Init()

	// 初始化虚拟机
	levm.CreateEVM()

	// 1 部署合约到state上
	contractAddress, _, err := levm.Deploy(caller, contractBinFile, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取调用的函数地址。函数名和参数名组成的字符串，哈希后的8个位
	input := crypto.Keccak256([]byte("get()"))[0:4]

	// 2 调用合约
	ret, _, err := levm.Call(caller, contractAddress, input, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(ret))
}
