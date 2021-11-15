# LonelyEVM
基于go-ethereum 1.9.25版本，针对EVM进行二次封装
## 快速开始
#### 安装
```
go get github.com/stevewooo/LonelyEVM@release-1.0
```
#### 导入
```golang
import (
  core "github.com/stevewooo/LonelyEVM/levm"
)
```
#### 初始化
```golang
var levm core.LEVM
levm.Init()
// 初始化虚拟机
levm.CreateEVM()
```
#### 合约部署
```golang
// 获得合约地址
contractAddress, _, err := levm.Deploy(caller, contractFile, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
```
#### 合约调用
```
input := crypto.Keccak256([]byte("get()"))[0:4] // 合约函数调用地址
ret, _, err := levm.Call(caller, contractAddress, input, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
fmt.Println(string(ret))
```
## 使用示范
本章介绍一个 Hello World 合约的编译与使用，此处使用的是0.7版本的solidity。
#### 编写合约
把以下内容保存到 ./contract.sol 文件中
```
pragma solidity ^0.7.0;
contract Hello {
    function get() public returns (string memory) {
        return "Hello World!";
    }
}
```
#### 安装solc与编译合约
1. 安装solcjs。
您也可以使用 http://remix.hubwiz.com/ 在线编译
```
npm install solc@0.7.4
```
2. 编译合约，获得编译后的 .bin 后缀文件，文件中的 hex 码内容就是我们需要传入到evm的内容。（remix 在线编译也能获得编译后的 hex 码内容）
```
node_modules/solc/solcjs --bin --abi contract.sol -o ./
# 得到 contract_sol_Hello.bin 文件和 contract_sol_Hello.abi 文件
```
#### 使用evm部署、调用合约
1. 在您的项目中导入本项目的核心，并从以太坊项目引入一些常用的包
```
import (
  core "github.com/stevewooo/LonelyEVM/levm"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/core/vm"
)
```
2. 然后在代码中对快速对虚拟机进行初始化
```
// 初始化
var levm core.LEVM
levm.Init()

// 创建虚拟机虚拟机
levm.CreateEVM()
```
3. 然后部署合约，并获得合约地址
```
// 初始化一些合约部署的内容（合约文件和调用者）
contractFile := common.Hex2Bytes("编译后的合约HEX码，作字符串传入即可")
caller := vm.AccountRef(common.HexToAddress("你的NODE_ID")) // 可以使用levm.EVMConfig.Origin代替

contractAddress, _, err := levm.Deploy(caller, contractFile, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
	if err != nil {
		fmt.Println(err)
		return
	}

```
4. 最后调用合约
```
// 获取调用的函数地址。函数名和参数名组成的字符串，哈希后的8个位
input := crypto.Keccak256([]byte("get()"))[0:4]

ret, _, err := levm.Call(caller, contractAddress, input, levm.EVMConfig.GasLimit, levm.EVMConfig.Value)
if err != nil {
  fmt.Println(err)
  return
}
fmt.Println(string(ret))
```
正常情况下能得到带空格的" Hello World"
