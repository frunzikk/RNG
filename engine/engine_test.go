package engine

import (
	"bytes"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"testing"
	"time"
)

func TestAccumulator(t *testing.T) {
	e := NewEngine()
	e.generator.reset()

	e.accumulator.addRandomEvent(0, 0, make([]byte, 32))
	e.accumulator.addRandomEvent(0, 0, make([]byte, 32))
	for i := uint(0); i < 1000; i++ {
		e.accumulator.addRandomEvent(1, i, []byte{1, 2})
	}
	out := e.RandomBytes(100)
	correct := []byte{
		159, 234, 36, 213, 49, 0, 53, 87, 54, 229, 161, 233, 221, 142, 54, 165, 76, 170, 203, 82, 54, 23, 41, 151, 61, 41, 163, 218, 50, 10, 160, 187, 135, 13, 208, 130, 125, 177, 215, 2, 53, 145, 151, 230, 24, 58, 189, 208, 117, 104, 111, 45, 119, 17, 166, 127, 103, 34, 91, 24, 113, 223, 19, 15, 52, 236, 212, 100, 242, 106, 9, 83, 110, 1, 169, 93, 31, 214, 97, 84, 80, 18, 220, 41, 105, 99, 164, 255, 148, 245, 186, 68, 106, 4, 202, 55, 42, 137, 137, 181,
	}
	if !bytes.Equal(out, correct) {
		t.Error("wrong RNG output")
		t.Error("generated", out)
		t.Error("correct", correct)
	}

	e.accumulator.addRandomEvent(0, 0, make([]byte, 32))
	e.accumulator.addRandomEvent(0, 0, make([]byte, 32))
	out = e.RandomBytes(100)
	correct = []byte{
		106, 105, 20, 236, 142, 230, 216, 120, 67, 101, 120, 123, 98, 73, 175, 227, 86, 151, 254, 60, 251, 139, 125, 63, 43, 192, 238, 86, 221, 105, 103, 225, 27, 87, 243, 92, 221, 67, 134, 104, 96, 79, 178, 4, 143, 108, 228, 194, 102, 205, 91, 225, 174, 24, 51, 227, 121, 95, 72, 190, 220, 195, 149, 153, 165, 244, 94, 112, 131, 85, 118, 14, 72, 109, 178, 84, 206, 80, 35, 31, 12, 29, 107, 51, 225, 34, 174, 115, 43, 181, 27, 136, 200, 226, 183, 22, 123, 167, 231, 193,
	}
	if !bytes.Equal(out, correct) {
		t.Error("wrong RNG output")
		t.Error("generated", out)
		t.Error("correct", correct)
	}

	time.Sleep(300 * time.Millisecond)

	out = e.RandomBytes(100)
	correct = []byte{
		12, 8, 249, 242, 85, 188, 147, 197, 122, 211, 43, 101, 20, 63, 111, 56, 119, 141, 105, 183, 23, 156, 128, 151, 68, 107, 66, 44, 196, 81, 68, 97, 38, 77, 70, 106, 149, 114, 10, 201, 167, 23, 227, 40, 174, 140, 38, 58, 135, 212, 17, 98, 25, 172, 82, 13, 144, 234, 246, 151, 227, 4, 146, 7, 211, 134, 182, 201, 41, 73, 224, 5, 40, 243, 239, 31, 1, 185, 192, 21, 242, 196, 98, 1, 88, 182, 104, 139, 103, 78, 175, 91, 103, 202, 88, 18, 122, 174, 11, 71,
	}
	if !bytes.Equal(out, correct) {
		t.Error("wrong RNG output")
		t.Error("generated", out)
		t.Error("correct", correct)
	}

}

func TestEngineInt63(t *testing.T) {
	e := NewEngine()
	e.Run()
	for i := 0; i < 100; i++ {
		r := e.Int63()
		if r < 0 {
			t.Error("Invalid random output")
		}
	}
}

func TestEngineUint64(t *testing.T) {
	e := NewEngine()
	e.Run()
	for i := 0; i < 100; i++ {
		r := e.Uint64()
		if r < 0 {
			t.Error("Invalid random output")
		}
	}
}

func TestEngineRand(t *testing.T) {
	e := NewEngine()
	e.Run()
	for i := 0; i < 100; i++ {
		r := e.GetRand(100, 0)
		if r < 0 {
			t.Error("Invalid random output")
		}
	}
}

// Benchmarking RNG
func engineRead(b *testing.B, n int) {
	e := NewEngine()
	e.Run()
	buffer := make([]byte, n)

	b.SetBytes(int64(n))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := e.Read(buffer); err != nil {
			b.Fatalf(err.Error())
		}
	}
}
func BenchmarkEngineRead16(b *testing.B) { engineRead(b, 16) }
func BenchmarkEngineRead32(b *testing.B) { engineRead(b, 32) }
func BenchmarkEngineRead1k(b *testing.B) { engineRead(b, 1024) }

// Benchmarking crypto/rand to compare with current implementation
func cryptoRandRead(b *testing.B, n int) {
	buffer := make([]byte, n)

	b.SetBytes(int64(n))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
			b.Fatalf(err.Error())
		}
	}
}
func BenchmarkCryptoRead16(b *testing.B) { cryptoRandRead(b, 16) }
func BenchmarkCryptoRead32(b *testing.B) { cryptoRandRead(b, 32) }
func BenchmarkCryptoRead1k(b *testing.B) { cryptoRandRead(b, 1024) }

func BenchmarkEngineInt63(b *testing.B) {
	e := NewEngine()
	e.Run()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Int63()
	}
}

func BenchmarkEngineUint64(b *testing.B) {
	e := NewEngine()
	e.Run()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Uint64()
	}
}

func BenchmarkMathInt63(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mrand.Int63()
	}
}

func BenchmarkMathUint64(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mrand.Uint64()
	}
}

func BenchmarkEngineRand(b *testing.B) {
	e := NewEngine()
	e.Run()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.GetRand(100, 0)
	}
}
