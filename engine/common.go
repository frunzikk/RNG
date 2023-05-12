package engine

import "encoding/binary"

func bytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

func uint64ToBytes(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, x)
	return bytes
}

func bytesToInt64(bytes []byte) int64 {
	return int64(bytesToUint64(bytes))
}

func int64ToBytes(x int64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(x))
	return bytes
}

func isZero(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

func wipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
