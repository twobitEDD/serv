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

// These accounts represent the affected accounts during the Claims decay bug

// The decay occurred before planned and the corresponding claimed amounts
// were less than supposed to be

package v12

// Accounts holds the missing claim amount to the corresponding account
var Accounts = [1][2]string{
	{"sx1000z4mk9jtr282yxmc3577qr9h6sprlu3tsz6l", "11874559096980682752"},
}
