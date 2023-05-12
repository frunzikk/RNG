package engine

import (
	"crypto/aes"
	"crypto/rand"
	"sync"
	"time"
)

type Engine struct {
	generator      *Generator
	generatorMutex sync.Mutex

	accumulator *Accumulator
}

func (engine *Engine) Run() {
	entropyChannel := engine.accumulator.EntropyDataChannel()
	go func() {
		source := engine.accumulator.allocateSource()
		for _ = range time.Tick(100 * time.Millisecond) {
			buffer := make([]byte, 4)
			rand.Read(buffer)
			entropyChannel <- struct {
				bytes  []byte
				source uint8
			}{bytes: buffer[len(buffer)-4:], source: source}
		}
	}()
	go func() {
		source := engine.accumulator.allocateSource()
		hasher := NewHasher()
		for _ = range time.Tick(100 * time.Millisecond) {
			buffer := make([]byte, 4)
			hasher.Write(int64ToBytes(time.Now().UnixNano()))
			buffer = append(buffer, hasher.Sum(nil)...)
			entropyChannel <- struct {
				bytes  []byte
				source uint8
			}{bytes: buffer[len(buffer)-4:], source: source}
		}
	}()
}

func (engine *Engine) RandomBytes(size uint) []byte {
	engine.generatorMutex.Lock()
	defer engine.generatorMutex.Unlock()
	engine.accumulator.checkReseed()
	return engine.generator.PseudoRandomData(size)
}

func (engine *Engine) RandomBytesUnlocked(size uint) []byte {
	engine.accumulator.checkReseed()
	return engine.generator.PseudoRandomData(size)
}

func (engine *Engine) GetRand(high uint64, low uint64) uint64 {
	bytes := engine.RandomBytes(8)
	bytes[0] &= 0x7f
	number := bytesToUint64(bytes)
	return (number % (high - low)) + low
}

func NewEngine() *Engine {
	generator := NewGenerator(aes.NewCipher)
	accumulator, _ := NewAccumulator(generator)
	engine := &Engine{
		generator:   generator,
		accumulator: accumulator,
	}
	return engine
}
