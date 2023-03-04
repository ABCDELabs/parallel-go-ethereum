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
	"sync"
	"testing"
	"time"
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
	pos2 := common.HexToHash("2")

	op1 := ResidualObject{Val: common.HexToHash("a"), Op: true}
	op2 := ResidualObject{Val: common.HexToHash("1"), Op: true}
	op3 := ResidualObject{Val: common.HexToHash("1"), Op: false}
	op4 := ResidualObject{Val: common.HexToHash("2"), Op: true}
	op5 := ResidualObject{Val: common.HexToHash("2"), Op: false}

	res := []ResidualObject{op1, op2}

	result := big.NewInt(0)
	for _, obj := range res {
		if obj.Op {
			result = result.Add(result, obj.Val.Big())
		} else {
			result = result.Sub(result, obj.Val.Big())
		}
	}

	expect := big.NewInt(11)
	if expect.Cmp(result) != 0 {
		t.Errorf("Merge Residual Object Test falied, expect %d, got %d", expect, result)
	} else {
		fmt.Printf("Merge Residual Object Test pass, expect %d, got %d\n", expect, result)
	}

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
	expect = big.NewInt(9)
	if expect.Cmp(result) != 0 {
		t.Errorf("Merge Residual Object Test falied, expect %d, got %d", expect, result)
	} else {
		fmt.Printf("Merge Residual Object Test pass, expect %d, got %d\n", expect, result)
	}

	state[pos2] = []ResidualObject{op1, op4, op5}
	result = big.NewInt(0)
	for _, obj := range state[pos2] {
		if obj.Op {
			result = result.Add(result, obj.Val.Big())
		} else {
			result = result.Sub(result, obj.Val.Big())
		}
	}
	expect = big.NewInt(10)
	if expect.Cmp(result) != 0 {
		t.Errorf("Merge Residual Object Test falied, expect %d, got %d", expect, result)
	} else {
		fmt.Printf("Merge Residual Object Test pass, expect %d, got %d\n", expect, result)
	}

}

func TestMockParallel(t *testing.T) {

	pos := common.HexToHash("1")
	val := common.HexToHash("2")
	//oppSub := common.HexToHash("0")
	oppADD := common.HexToHash("1")

	tsobj := NewTStateObject()

	tsobj.dirtyStorage[pos] = val

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			go tsobj.SetResidualState(pos, val, oppADD)
		} else {
			go tsobj.SetResidualState(pos, val, oppADD)
		}
	}

	time.Sleep(time.Second * 2)

	tsobj.MergeResidualState()

	fmt.Println(tsobj.dirtyStorage[pos].Big())

}

func TestMockConcurrentTransaction(t *testing.T) {
	var wg sync.WaitGroup

	tsobj := NewTStateObject()

	pos := common.HexToHash("1")
	val := common.HexToHash("2")
	oppADD := common.HexToHash("1")

	start := time.Now()

	for i := 0; i < 1000; i++ {
		val = common.BigToHash(val.Big().Add(val.Big(), big.NewInt(1)))
		time.Sleep(time.Millisecond * 8)
		tsobj.SetState(pos, val)
	}
	elapsed := time.Since(start)
	fmt.Println("The sequential execution result:", tsobj.dirtyStorage[pos].Big())
	fmt.Println("The sequential execution time:", elapsed)

	tsobj.dirtyStorage = make(Storage)

	tsobj.dirtyStorage[pos] = common.HexToHash("2")
	val = common.HexToHash("1")

	// Parallel Mode
	start = time.Now()
	// Seq num
	for i := 0; i < 901; i++ {
		val = common.BigToHash(val.Big().Add(val.Big(), big.NewInt(1)))
		time.Sleep(time.Millisecond * 8)
		tsobj.SetState(pos, val)
	}
	val = common.HexToHash("1")
	// Parallel Number
	for i := 0; i < 100; i++ {
		wg.Add(1)
		// Parallel
		go func() {
			tsobj.SetResidualState(pos, val, oppADD)
			time.Sleep(time.Millisecond * 8)
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed = time.Since(start)
	tsobj.MergeResidualState()
	fmt.Println("The Parallel execution result:", tsobj.dirtyStorage[pos].Big())
	fmt.Println("The Parallel execution time:", elapsed)

}
