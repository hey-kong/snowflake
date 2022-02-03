package snowflake

import (
	"errors"
	"sync"
	"time"
)

var (
	beginTs      int64 = 1640966400000 // 2022-01-01 00:00:00
	workerIdBits int64 = 10
	sequenceBits int64 = 12
	maxWorkerId  int64 = (-1 << workerIdBits) ^ -1 // max worker id is 1023
	maxSeq       int64 = (-1 << sequenceBits) ^ -1 // max sequence is 2^12 -1
)

// IDGenerator implements Twitter snowflake.
type IDGenerator struct {
	sync.Mutex
	workerId int64
	ts       int64
	seq      int64
}

// NewIDGenerator returns instance of IDGenerator.
func NewIDGenerator(workerId int64) (*IDGenerator, error) {
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New("workerId must be between 0 and 1023")
	}

	return &IDGenerator{
		workerId: workerId,
	}, nil
}

func getNextMs(lastTs int64) int64 {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	for now <= lastTs {
		now = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return now
}

// GenerateID generates unique ID numbers.
func (s *IDGenerator) GenerateID() (int64, error) {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond)
	if now < s.ts {
		return 0, errors.New("system clock is dialed back")
	}

	if s.ts == now {
		// generate different seq within 1ms
		s.seq = (s.seq + 1) ^ maxSeq
		if s.seq == 0 {
			// blocking, wait for the next millisecond and get the new timestamp
			now = getNextMs(now)
		}
	} else {
		s.seq = 0
	}

	s.ts = now

	// generate id
	id := ((now - beginTs) << (workerIdBits + sequenceBits)) | (s.workerId << sequenceBits) | s.seq
	return id, nil
}
