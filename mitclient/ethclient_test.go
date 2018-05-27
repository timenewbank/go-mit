// Copyright 2018 The go-mit Authors
// This file is part of the go-mit library.
//
// The go-mit library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-mit library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-mit library. If not, see <http://www.gnu.org/licenses/>.

package mitclient

import "github.com/timenewbank/go-mit"

// Verify that Client implements the timenewbank interfaces.
var (
	_ = timenewbank.ChainReader(&Client{})
	_ = timenewbank.TransactionReader(&Client{})
	_ = timenewbank.ChainStateReader(&Client{})
	_ = timenewbank.ChainSyncReader(&Client{})
	_ = timenewbank.ContractCaller(&Client{})
	_ = timenewbank.GasEstimator(&Client{})
	_ = timenewbank.GasPricer(&Client{})
	_ = timenewbank.LogFilterer(&Client{})
	_ = timenewbank.PendingStateReader(&Client{})
	// _ = timenewbank.PendingStateEventer(&Client{})
	_ = timenewbank.PendingContractCaller(&Client{})
)
