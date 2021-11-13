package main

import (
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// Lonely EVM，封装了调用evm的方法。
type LEVM struct {
	UserConfig UserConfig
	EVMConfig  EVMConfig
	Evm        *vm.EVM
}

func (levm *LEVM) Init() {
	// 初始化evm配置
	var cfg EVMConfig
	cfg.SetDefaults()

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	}

	// 创建一台EVM
	txContext := vm.TxContext{
		Origin:   cfg.Origin,
		GasPrice: cfg.GasPrice,
	}
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     cfg.GetHashFn,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
	}
	levm.Evm = vm.NewEVM(blockContext, txContext, cfg.State, cfg.ChainConfig, cfg.EVMConfig)

	if cfg.ChainConfig.IsYoloV2(levm.Evm.Context.BlockNumber) {
		// cfg.State.AddAddressToAccessList(cfg.Origin) // 添加自己的NodeID为白名单
		for _, addr := range levm.Evm.ActivePrecompiles() {
			cfg.State.AddAddressToAccessList(addr)
		}
	}

	levm.EVMConfig = cfg
}

// 部署合约
func (levm *LEVM) Deploy(caller vm.ContractRef, code []byte, gas uint64, value *big.Int) (common.Address, uint64, error) {
	_, contractAddress, leftOverGas, err := levm.Evm.Create(
		caller,
		code,
		gas,
		value,
	)

	return contractAddress, leftOverGas, err
}

// 调用合约
func (levm *LEVM) Call(caller vm.ContractRef, contractAddress common.Address, input []byte, gas uint64, value *big.Int) ([]byte, uint64, error) {
	// 调用合约
	return levm.Evm.Call(
		caller,
		contractAddress,
		input,
		gas,
		value,
	)
}

type EVMConfig struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config
	BaseFee     *big.Int

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash
}

// sets defaults on the config
func (cfg *EVMConfig) SetDefaults() {
	if cfg.ChainConfig == nil {
		cfg.ChainConfig = &params.ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      new(big.Int),
			DAOForkBlock:        new(big.Int),
			DAOForkSupport:      false,
			EIP150Block:         new(big.Int),
			EIP150Hash:          common.Hash{},
			EIP155Block:         new(big.Int),
			EIP158Block:         new(big.Int),
			ByzantiumBlock:      new(big.Int),
			ConstantinopleBlock: new(big.Int),
			PetersburgBlock:     new(big.Int),
			IstanbulBlock:       new(big.Int),
			MuirGlacierBlock:    new(big.Int),
			YoloV2Block:         nil,
		}
	}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
}
