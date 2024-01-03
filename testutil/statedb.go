// Copyright 2022 Serv Foundation
// This file is part of the Serv Network packages.
//
// Serv is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Serv packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Serv packages. If not, see https://github.com/twobitedd/serv/blob/main/LICENSE

package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/twobitedd/serv/v12/app/ante/evm"
	"github.com/twobitedd/serv/v12/x/evm/statedb"
)

// NewStateDB returns a new StateDB for testing purposes.
func NewStateDB(ctx sdk.Context, evmKeeper evm.EVMKeeper) *statedb.StateDB {
	return statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes())))
}
