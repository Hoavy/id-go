// Copyright 2017 wanghaowei <d@wanghaowei.com> All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package idflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	epoch              = uint64(1488888888123)
	workerIdBits       = uint(10)
	sequenceBits       = uint(12)
	workerIdShift      = sequenceBits
	timestampLeftShift = sequenceBits + workerIdBits
	maxWorkerId        = -1 ^ (-1 << workerIdBits)
	maxSequenceId      = -1 ^ (-1 << sequenceBits)
	maxTimestamp       = -1 ^ (-1 << (64 - timestampLeftShift))
	maxId              = -1 ^ (-1 << 64)
)

type Idflake struct {
	epoch         uint64
	sequence      uint64
	lastTimestamp uint64
	workerId      uint64
	mutex         *sync.Mutex
}

func NewIdflake(workerId uint64) (*Idflake, error) {
	idflake := &Idflake{}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("worker Id: %d error", workerId))
	}

	idflake.epoch = epoch
	idflake.sequence = 0
	idflake.lastTimestamp = 1
	idflake.workerId = workerId
	idflake.mutex = &sync.Mutex{}
	return idflake, nil
}

func (id *Idflake) NextId() (uint64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	timestamp := id.timeGen()

	if (timestamp - id.epoch) >= maxTimestamp {
		return 0, errors.New(fmt.Sprintf("Timestamp overflows %d milliseconds", maxTimestamp))
	}

	if timestamp < id.lastTimestamp {
		return 0, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp))
	}

	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1)
		if (id.sequence ^ maxSequenceId) == 0 {
			timestamp = id.skipNextMillis(id.lastTimestamp)
			id.sequence = 0
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp

	return (timestamp-id.epoch)<<timestampLeftShift | id.workerId<<workerIdShift | id.sequence, nil
}

func (id *Idflake) SetEpoch(epoch uint64) (bool, error) {
	id.epoch = epoch
	return true, nil
}

func (id *Idflake) skipNextMillis(lastTimestamp uint64) uint64 {
	timestamp := id.timeGen()
	for timestamp <= lastTimestamp {
		timestamp = id.timeGen()
	}
	return timestamp
}

func (id *Idflake) timeGen() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}
