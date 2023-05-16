package engine

import (
	"hash"
	"sync"
	"time"
)

const (
	poolCount         = 32
	minPoolSize       = 32
	reseedInterval    = 100 * time.Millisecond
	channelBufferSize = 4
)

type Accumulator struct {
	generator     *Generator
	pool          [poolCount]hash.Hash
	poolMutex     sync.Mutex
	reseedCount   int
	nextReseed    time.Time
	firstPoolSize int
	sourceMutex   sync.Mutex
	nextSource    uint8
	stopSources   chan bool
	sources       sync.WaitGroup
}

func (accumulator *Accumulator) checkReseed() {
	now := time.Now()

	accumulator.poolMutex.Lock()
	defer accumulator.poolMutex.Unlock()

	if accumulator.firstPoolSize >= minPoolSize && now.After(accumulator.nextReseed) {
		accumulator.nextReseed = now.Add(reseedInterval)
		accumulator.firstPoolSize = 0
		accumulator.reseedCount++

		seed := make([]byte, 0, poolCount*HasherSize)
		for i := 0; i < poolCount; i++ {
			mask := 1 << i
			if accumulator.reseedCount%mask != 0 {
				break
			}
			seed = accumulator.pool[i].Sum(seed)
			accumulator.pool[i].Reset()
		}
		accumulator.generator.reseed(seed)
	}
}

func (accumulator *Accumulator) allocateSource() uint8 {
	accumulator.sourceMutex.Lock()
	defer accumulator.sourceMutex.Unlock()
	source := accumulator.nextSource
	accumulator.nextSource++
	accumulator.sources.Add(1)
	return source
}

func (accumulator *Accumulator) entropyDataChannel() chan<- struct {
	bytes  []byte
	source uint8
} {
	c := make(chan struct {
		bytes  []byte
		source uint8
	}, channelBufferSize)

	go func() {
		defer accumulator.sources.Done()
		seq := uint(0)

	loop:
		for {
			select {
			case data, ok := <-c:
				if !ok {
					break loop
				}

				if len(data.bytes) > 32 {
					hasher := NewHasher()
					hasher.Write(data.bytes)
					data.bytes = hasher.Sum(nil)
				}

				accumulator.addRandomEvent(data.source, seq, data.bytes)
				seq++
			case <-accumulator.stopSources:
				break loop
			}
		}
	}()

	return c
}

func (accumulator *Accumulator) addRandomEvent(source uint8, seq uint, data []byte) {
	poolNumber := seq % poolCount
	accumulator.poolMutex.Lock()

	pool := accumulator.pool[poolNumber]
	pool.Write([]byte{source, byte(len(data))})
	pool.Write(data)
	if poolNumber == 0 {
		accumulator.firstPoolSize += 2 + len(data)
	}

	accumulator.poolMutex.Unlock()
	accumulator.checkReseed()
}

func NewAccumulator(generator *Generator) (*Accumulator, error) {
	accumulator := &Accumulator{generator: generator}
	for i := 0; i < len(accumulator.pool); i++ {
		accumulator.pool[i] = NewHasher()
	}
	accumulator.stopSources = make(chan bool)
	return accumulator, nil
}
