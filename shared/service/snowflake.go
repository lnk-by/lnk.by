package service

import (
	"errors"
	"hash/fnv"
	"os"
	"sync"
	"time"
)

const (
	epoch         int64 = 1577836800000 // Custom epoch: 2020-01-01T00:00:00Z
	timestampBits       = 41
	machineIDBits       = 10
	sequenceBits        = 12

	maxMachineID   = -1 ^ (-1 << machineIDBits) // 1023
	maxSequence    = -1 ^ (-1 << sequenceBits)  // 4095
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type Generator struct {
	mu         sync.Mutex
	lastUnixMs int64
	sequence   int64
	machineID  int64
}

func NewDefaultGenerator() (*Generator, error) {
	return NewGenerator(GetMachineID())
}

// NewGenerator creates a new Snowflake ID generator with the given machine ID (0–1023).
func NewGenerator(machineID int64) (*Generator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID out of range")
	}
	return &Generator{
		machineID: machineID,
	}, nil
}

// NextID generates a new unique 64-bit Snowflake ID.
func (g *Generator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastUnixMs {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// Sequence overflow, wait for next millisecond
			for now <= g.lastUnixMs {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastUnixMs = now

	return ((now - epoch) << timestampShift) |
		(g.machineID << machineIDShift) |
		g.sequence
}

func GetMachineID() int64 {
	val := os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME")
	if val == "" {
		val = "default-machine"
	}
	h := fnv.New32a()
	h.Write([]byte(val))
	return int64(h.Sum32() & 0x3FF) // 10 bits = 1023 max
}

func EncodeBase62(n int64) string {
	if n == 0 {
		return string(base62Chars[0])
	}
	result := make([]byte, 0)
	for n > 0 {
		result = append([]byte{base62Chars[n%62]}, result...)
		n /= 62
	}
	return string(result)
}
