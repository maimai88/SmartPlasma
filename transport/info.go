package transport

import (
	"math/big"

	"github.com/SmartMeshFoundation/Spectrum/common"
	"github.com/pkg/errors"

	"github.com/SmartMeshFoundation/SmartPlasma/contract/rootchain"
	"github.com/SmartMeshFoundation/SmartPlasma/transport/handlers"
)

// DepositCount returns a deposit counter.
func (c *Client) DepositCount() (count *big.Int, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		return session.DepositCount()
	}
	req := &handlers.DepositCountReq{}
	var resp *handlers.DepositCountResp
	call := c.connect.Go(DepositCountMethod, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return nil, replay.Error
		}
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp.Count, err
}

// ChallengePeriod returns a period for challenging in seconds.
func (c *Client) ChallengePeriod() (count *big.Int, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		return session.ChallengePeriod()
	}

	req := &handlers.ChallengePeriodReq{}
	var resp *handlers.ChallengePeriodResp
	call := c.connect.Go(ChallengePeriodMethod, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return nil, replay.Error
		}
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp.ChallengePeriod, err
}

// Operator returns a Plasma Cash operator address.
func (c *Client) Operator() (address common.Address, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		return session.Operator()
	}
	req := &handlers.OperatorReq{}
	var resp *handlers.OperatorResp
	call := c.connect.Go(OperatorMethod, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return common.Address{}, replay.Error
		}
	case <-ctx.Done():
		return common.Address{}, errors.New("timeout")
	}

	if resp.Error != "" {
		return common.Address{}, errors.New(resp.Error)
	}

	return resp.Operator, err
}

// ChildChain returns a block hash by a block number.
func (c *Client) ChildChain(
	blockNumber *big.Int) (hash common.Hash, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		return session.ChildChain(blockNumber)
	}
	req := &handlers.ChildChainReq{
		BlockNumber: blockNumber,
	}
	var resp *handlers.ChildChainResp
	call := c.connect.Go(ChildChainMethod, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return common.Hash{}, replay.Error
		}
	case <-ctx.Done():
		return common.Hash{}, errors.New("timeout")
	}

	if resp.Error != "" {
		return common.Hash{}, errors.New(resp.Error)
	}

	return resp.BlockHash, err
}

// Exits returns a incomplete exit by UID.
func (c *Client) Exits(uid *big.Int) (resp *handlers.ExitsResp, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		result, err := session.Exits(uid)
		if err != nil {
			return nil, err
		}
		resp = &handlers.ExitsResp{
			State:                result.State,
			ExitTime:             result.ExitTime,
			ExitTxBlkNum:         result.ExitTxBlkNum,
			ExitTx:               result.ExitTx,
			TxBeforeExitTxBlkNum: result.TxBeforeExitTxBlkNum,
			TxBeforeExitTx:       result.TxBeforeExitTx,
		}
		return resp, err
	}
	req := &handlers.ExitsReq{
		UID: uid,
	}

	call := c.connect.Go(ExitsMethod, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return nil, replay.Error
		}
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp, err
}

// Wallet returns a deposit amount.
func (c *Client) Wallet(uid *big.Int) (amount *big.Int, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		return session.Wallet(common.BigToHash(uid))
	}
	req := &handlers.WalletReq{
		UID: uid,
	}
	var resp *handlers.WalletResp
	call := c.connect.Go(WalletMethod, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return nil, replay.Error
		}
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp.Amount, err
}

// Wallet2 returns a deposit Smart Plasma block number.
func (c *Client) Wallet2(uid *big.Int) (block *big.Int, err error) {
	ctx, cancel := c.newContext()
	defer cancel()

	if c.sessionRootChain != nil {
		session := rootchain.CopySession(c.sessionRootChain)
		session.TransactOpts.Context = ctx
		return session.Wallet2(uid)
	}
	req := &handlers.Wallet2Req{
		UID: uid,
	}
	var resp *handlers.Wallet2Resp
	call := c.connect.Go(Wallet2Method, req, &resp, nil)

	select {
	case replay := <-call.Done:
		if replay.Error != nil {
			return nil, replay.Error
		}
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp.BlockNumber, err
}
