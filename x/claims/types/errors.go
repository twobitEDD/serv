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

package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrClaimsRecordNotFound = errorsmod.Register(ModuleName, 2, "claims record not found")
	ErrInvalidAction        = errorsmod.Register(ModuleName, 3, "invalid claim action type")
	ErrKeyTypeNotSupported  = errorsmod.Register(ModuleName, 4, "key type 'secp256k1' not supported")
)
