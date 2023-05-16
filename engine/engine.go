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
	entropyChannel := engine.accumulator.entropyDataChannel()
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

func (engine *Engine) Read(p []byte) (n int, err error) {
	copy(p, engine.RandomBytes(uint(len(p))))
	return len(p), nil
}

func (engine *Engine) RandomBytes(size uint) []byte {
	engine.generatorMutex.Lock()
	defer engine.generatorMutex.Unlock()
	engine.accumulator.checkReseed()
	return engine.generator.pseudoRandomData(size)
}

func (engine *Engine) RandomBytesUnlocked(size uint) []byte {
	engine.accumulator.checkReseed()
	return engine.generator.pseudoRandomData(size)
}

func (engine *Engine) GetRand(high int64, low int64) int64 {
	bytes := engine.RandomBytes(8)
	bytes[0] &= 0x7f
	number := bytesToInt64(bytes)
	return (number % (high - low)) + low
}

func (engine *Engine) Int63() int64 {
	engine.generatorMutex.Lock()
	defer engine.generatorMutex.Unlock()
	randomBytes := engine.generator.pseudoRandomData(8)
	randomBytes[0] &= 0x7f
	return bytesToInt64(randomBytes)
}

func (engine *Engine) Seed(seed int64) {
	engine.generatorMutex.Lock()
	defer engine.generatorMutex.Unlock()
	engine.generator.reset()
	engine.generator.reseedInt64(seed)
}

func (engine *Engine) Uint64() uint64 {
	randomBytes := engine.generator.pseudoRandomData(8)
	return bytesToUint64(randomBytes)
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
