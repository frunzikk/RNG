package engine

import (
	"crypto/sha256"
	"hash"
)

const (
	HasherSize = 32
	BlockSize  = 64
)

type Hasher struct {
	firstRound  hash.Hash
	secondRound hash.Hash
}

func (hash *Hasher) Reset() {
	hash.firstRound.Reset()
}

func (hash *Hasher) Write(p []byte) (n int, err error) {
	return hash.firstRound.Write(p)
}

func (hash *Hasher) Sum(b []byte) []byte {
	hash.secondRound.Reset()
	hash.secondRound.Write(hash.firstRound.Sum(nil))
	return hash.secondRound.Sum(b)
}

func (hash *Hasher) Size() int {
	return HasherSize
}

func (hash *Hasher) BlockSize() int {
	return BlockSize
}

func NewHasher() hash.Hash {
	return &Hasher{firstRound: sha256.New(), secondRound: sha256.New()}
}
