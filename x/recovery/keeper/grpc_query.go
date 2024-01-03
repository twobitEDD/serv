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

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/twobitedd/serv/v12/x/recovery/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(
	c context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}
