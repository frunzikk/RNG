package engine

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"net"
	"os/user"
	"time"
)

const (
	keySize   = HasherSize
	maxBlocks = 1 << 16
)

type NewCipher func([]byte) (cipher.Block, error)

type Generator struct {
	newCipher NewCipher
	key       []byte
	cipher    cipher.Block
	counter   []byte
}

func (generator *Generator) increment() {
	counter := generator.counter
	for i := 0; i < len(counter); i++ {
		counter[i]++
		if counter[i] != 0 {
			break
		}
	}
}

func (generator *Generator) reseed(seed []byte) {
	hash := NewHasher()
	hash.Write(generator.key)
	hash.Write(seed)
	generator.setKey(hash.Sum(nil))
	generator.increment()
}

func (generator *Generator) setKey(key []byte) {
	if len(key) != keySize {
		panic("Wrong key size")
	}
	generator.key = key
	newCipher, err := generator.newCipher(generator.key)
	if err != nil {
		panic("cannot set generator key")
	}
	generator.cipher = newCipher
}

func (generator *Generator) reset() {
	zeroKey := make([]byte, keySize)
	generator.setKey(zeroKey)
	generator.counter = make([]byte, generator.cipher.BlockSize())
}

func (generator *Generator) generateBlocks(data []byte, size uint) []byte {
	buffer := make([]byte, len(generator.counter))
	for i := uint(0); i < size; i++ {
		generator.cipher.Encrypt(buffer, generator.counter)
		data = append(data, buffer...)
		generator.increment()
	}
	return data
}

func (generator *Generator) setInitialSeed() {
	seedData := &bytes.Buffer{}
	isGood := false

	written, _ := io.CopyN(seedData, rand.Reader, keySize)
	isGood = isGood || (written >= keySize)

	if !isGood {
		panic("failed to get initial randomness for the seed")
	}

	// current time (different between different runs of the program, difficult to predict)
	now := time.Now()
	seedData.Write(int64ToBytes(now.UnixNano()))

	// network interfaces (different between hosts)
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		seedData.Write(int64ToBytes(int64(iface.MTU)))
		seedData.Write([]byte(iface.Name))
		seedData.Write(iface.HardwareAddr)
		seedData.Write(int64ToBytes(int64(iface.Flags)))
	}

	// user account details (maybe different between hosts)
	currentUser, _ := user.Current()
	if currentUser != nil {
		seedData.Write([]byte(currentUser.Uid))
		seedData.Write([]byte(currentUser.Gid))
		seedData.Write([]byte(currentUser.Username))
		seedData.Write([]byte(currentUser.Name))
		seedData.Write([]byte(currentUser.HomeDir))
	}

	buf := seedData.Bytes()
	generator.reseed(buf)
	wipe(buf)
}

func (generator *Generator) numBlocks(n uint) uint {
	k := uint(len(generator.counter))
	return (n + k - 1) / k
}

func (generator *Generator) pseudoRandomData(size uint) []byte {
	numBlocks := generator.numBlocks(size)
	res := make([]byte, 0, numBlocks*uint(len(generator.counter)))

	for numBlocks > 0 {
		count := numBlocks
		if count > numBlocks {
			count = maxBlocks
		}
		res = generator.generateBlocks(res, count)
		numBlocks -= count

		newKey := generator.generateBlocks(nil, generator.numBlocks(keySize))
		generator.setKey(newKey[:keySize])
	}

	return res[:size]
}

func (generator *Generator) reseedInt64(seed int64) {
	seedBytes := int64ToBytes(seed)
	generator.reseed(seedBytes)
}

func NewGenerator(newCipher NewCipher) *Generator {
	generator := &Generator{
		newCipher: newCipher,
	}
	generator.reset()
	generator.setInitialSeed()
	return generator
}
