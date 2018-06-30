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

package params

import "math/big"

var (
	TargetGasLimit uint64 = GenesisGasLimit // The artificial target
)

const (
	GasLimitBoundDivisor uint64 = 2048    // 1024The bound divisor of the gas limit, used in update calculations.
	MinGasLimit          uint64 = 10000    // 5000Minimum the gas limit may ever be.
	GenesisGasLimit      uint64 = 9424776 // 4712388Gas limit of the Genesis block.

	MaximumExtraDataSize  uint64 = 32    // Maximum size extra data may be after Genesis.
	ExpByteGas            uint64 = 10    // Times ceil(log256(exponent)) for the EXP instruction.
	SloadGas              uint64 = 50    // Multiplied by the number of 32-byte words that are copied (round up) for any *COPY operation and added.
	CallValueTransferGas  uint64 = 18000  // 9000Paid for CALL when the value transfer is non-zero.
	CallNewAccountGas     uint64 = 50000 // 25000Paid for CALL when the destination address didn't exist prior.
	TxGas                 uint64 = 42000 // 21000Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	TxGasContractCreation uint64 = 106000 // 53000Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.
	TxDataZeroGas         uint64 = 8     // 4Per byte of data attached to a transaction that equals zero. NOTE: Not payable on data of calls between transactions.
	QuadCoeffDiv          uint64 = 512   // Divisor for the quadratic particle of the memory cost equation.
	SstoreSetGas          uint64 = 40000 // 20000Once per SLOAD operation.
	LogDataGas            uint64 = 16     // 8Per byte in a LOG* operation's data.
	CallStipend           uint64 = 2300  // Free gas given at beginning of call.

	Sha3Gas          uint64 = 60    // 30Once per SHA3 operation.
	Sha3WordGas      uint64 = 12     // 6Once per word of the SHA3 operation's data.
	SstoreResetGas   uint64 = 10000  // 5000Once per SSTORE operation if the zeroness changes from zero.
	SstoreClearGas   uint64 = 10000  // 5000Once per SSTORE operation if the zeroness doesn't change.
	SstoreRefundGas  uint64 = 30000 // 15000Once per SSTORE operation if the zeroness changes to zero.
	JumpdestGas      uint64 = 2     // 1Refunded gas, once per SSTORE operation if the zeroness changes to zero.
	EpochDuration    uint64 = 30000 // Duration between proof-of-work epochs.
	CallGas          uint64 = 40    // Once per CALL operation & message call transaction.
	CreateDataGas    uint64 = 400   // 200
	CallCreateDepth  uint64 = 1024  // Maximum depth of call/create stack.
	ExpGas           uint64 = 10    // Once per EXP instruction
	LogGas           uint64 = 750   // 375Per LOG* operation.
	CopyGas          uint64 = 6     // 3
	StackLimit       uint64 = 1024  // Maximum size of VM stack allowed.
	TierStepGas      uint64 = 0     // Once per operation, for a selection of them.
	LogTopicGas      uint64 = 750   // 375Multiplied by the * of the LOG*, per LOG transaction. e.g. LOG0 incurs 0 * c_txLogTopicGas, LOG4 incurs 4 * c_txLogTopicGas.
	CreateGas        uint64 = 64000 // 32000Once per CREATE operation & contract-creation transaction.
	SuicideRefundGas uint64 = 48000 // 24000Refunded following a suicide operation.
	MemoryGas        uint64 = 6     // 3Times the address of the (highest referenced byte in memory + 1). NOTE: referencing happens on read, write and in instructions such as RETURN and CALL.
	TxDataNonZeroGas uint64 = 136    // 68Per byte of data attached to a transaction that is not equal to zero. NOTE: Not payable on data of calls between transactions.

	MaxCodeSize = 24576 // Maximum bytecode to permit for a contract

	// Precompiled contract gas prices

	EcrecoverGas            uint64 = 6000   // 3000Elliptic curve sender recovery gas price
	Sha256BaseGas           uint64 = 120     // 60Base price for a SHA256 operation
	Sha256PerWordGas        uint64 = 24     // 12Per-word price for a SHA256 operation
	Ripemd160BaseGas        uint64 = 1200    // 600Base price for a RIPEMD160 operation
	Ripemd160PerWordGas     uint64 = 240    // 120Per-word price for a RIPEMD160 operation
	IdentityBaseGas         uint64 = 30     // 15Base price for a data copy operation
	IdentityPerWordGas      uint64 = 6      // 3Per-work price for a data copy operation
	ModExpQuadCoeffDiv      uint64 = 20     // Divisor for the quadratic particle of the big int modular exponentiation
	Bn256AddGas             uint64 = 1000    // 500Gas needed for an elliptic curve addition
	Bn256ScalarMulGas       uint64 = 80000  // 40000Gas needed for an elliptic curve scalar multiplication
	Bn256PairingBaseGas     uint64 = 200000 // 100000Base price for an elliptic curve pairing check
	Bn256PairingPerPointGas uint64 = 160000  // 80000Per-point price for an elliptic curve pairing check
)

var (
	DifficultyBoundDivisor = big.NewInt(2048)   // The bound divisor of the difficulty, used in the update calculations.
	GenesisDifficulty      = big.NewInt(131072) // Difficulty of the Genesis block.
	MinimumDifficulty      = big.NewInt(131072) // The minimum that the difficulty may ever be.
	DurationLimit          = big.NewInt(13)     // The decision boundary on the blocktime duration used to determine whether difficulty should go up or not.
)
