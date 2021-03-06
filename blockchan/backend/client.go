package backend

import (
	"context"
	"log"
	"time"

	"github.com/SmartMeshFoundation/Spectrum/accounts/abi/bind"
	"github.com/SmartMeshFoundation/Spectrum/accounts/abi/bind/backends"
	"github.com/SmartMeshFoundation/Spectrum/common"
	"github.com/SmartMeshFoundation/Spectrum/core/types"
	"github.com/SmartMeshFoundation/Spectrum/ethclient"
	"github.com/pkg/errors"
)

// Errors.
var (
	ErrInvalidBackend = errors.New("unsupported backend")
)

// Backend interface.
type Backend interface {
	Connect() bind.ContractBackend
	Mine(ctx context.Context, tx *types.Transaction) (*types.Receipt, error)
	GoodTransaction(tx *types.Transaction) bool
}

// Simulator interface.
type Simulator interface {
	Backend
	AdjustTime(adjustment time.Duration) error
}

type backend struct {
	connect bind.ContractBackend
}

// NewBackend makes new Backend.
func NewBackend(url string) Backend {
	cli, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &backend{connect: cli}
}

// NewSimulatedBackend makes new backend simulator.
func NewSimulatedBackend(accounts []common.Address) Backend {
	return &backend{connect: newSimulator(accounts)}
}

// Connect gets connect to Ethereum backend.
func (back *backend) Connect() bind.ContractBackend {
	return back.connect
}

// Mine to wait mining.
func (back *backend) Mine(ctx context.Context,
	tx *types.Transaction) (*types.Receipt, error) {
	switch conn := back.connect.(type) {
	case *ethclient.Client:
		return bind.WaitMined(ctx, conn, tx)
	case *backends.SimulatedBackend:
		conn.Commit()
		return bind.WaitMined(ctx, conn, tx)
	}
	return nil, ErrInvalidBackend
}

// AdjustTime adds a time shift to the simulated clock.
func (back *backend) AdjustTime(adjustment time.Duration) error {
	switch conn := back.connect.(type) {
	case *backends.SimulatedBackend:
		err := conn.AdjustTime(adjustment)
		if err != nil {
			return err
		}
		conn.Commit()
	}
	return nil
}

// GoodTransaction returns true if transaction status = 1.
func (back *backend) GoodTransaction(tx *types.Transaction) bool {
	tr, err := back.Mine(context.Background(), tx)
	if err != nil {
		return false
	}

	//fmt.Printf("gas: %d\n", tr.GasUsed)

	if tr.Status != 1 {
		return false
	}
	return true
}
