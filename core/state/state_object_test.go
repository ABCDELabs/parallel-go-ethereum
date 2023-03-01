// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"bytes"
	"fmt"
	"github.com/ABCDEcapital/parallel-go-ethereum/common"
	"math/big"
	"testing"
)

func BenchmarkCutOriginal(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		bytes.TrimLeft(value[:], "\x00")
	}
}

func BenchmarkCutsetterFn(b *testing.B) {
	value := common.HexToHash("0x01")
	cutSetFn := func(r rune) bool { return r == 0 }
	for i := 0; i < b.N; i++ {
		bytes.TrimLeftFunc(value[:], cutSetFn)
	}
}

func BenchmarkCutCustomTrim(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		common.TrimLeftZeroes(value[:])
	}
}

// TODO: Dismiss
func TestMergeResidualState(t *testing.T) {
	s := common.HexToHash("a")
	result := big.NewInt(0)

	vm := common.BytesToHash([]byte("0xfffffffffffffffffffffffffffffffd"))
	testValues := []common.Hash{common.HexToHash("a"), vm}
	for _, v := range testValues {
		fmt.Println("The value of v:", v)
	}

	fmt.Println("The value of Hash:", s.Big())
	fmt.Println("The value of Merge:", result)

	d := big.NewInt(-3)
	if d.Cmp(big.NewInt(0)) >= 0 {
		fmt.Println("The result is great than 0", common.BigToHash(big.NewInt(-3)).Big())
	} else {
		fmt.Println("The result is less than 0", common.BigToHash(big.NewInt(-3)))
	}
}

func TestMergeResidualObject(t *testing.T) {

	pos1 := common.HexToHash("1")

	op1 := ResidualObject{Val: common.HexToHash("a"), Op: true}
	op2 := ResidualObject{Val: common.HexToHash("1"), Op: true}
	op3 := ResidualObject{Val: common.HexToHash("2"), Op: false}

	res := []ResidualObject{op1, op2}

	fmt.Println(op1.Val.Big(), op2.Val.Big())

	result := big.NewInt(0)
	for _, obj := range res {
		if obj.Op {
			result = result.Add(result, obj.Val.Big())
		} else {
			result = result.Sub(result, obj.Val.Big())
		}
	}
	fmt.Println("The Merge Result", result)

	state := make(ResidualStorage)
	state[pos1] = []ResidualObject{op1, op3}
	result = big.NewInt(0)
	for _, obj := range state[pos1] {
		if obj.Op {
			result = result.Add(result, obj.Val.Big())
		} else {
			result = result.Sub(result, obj.Val.Big())
		}
	}
	fmt.Println("The Merge Result", result)

}
