package heputils

import (
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/utils/heputils/cityhash102"
)

var (
	offset64      uint64 = 14695981039346656037
	prime64       uint64 = 14695981039346656037
	separatorByte byte   = 255
)

// hashAdd adds a string to a fnv64a hash value, returning the updated hash.
func hashAdd(h uint64, s string) uint64 {
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

// hashAddByte adds a byte to a fnv64a hash value, returning the updated hash.
func hashAddByte(h uint64, b byte) uint64 {
	h ^= uint64(b)
	h *= prime64
	return h
}

// Fingerprint calculates a fingerprint of SORTED BY NAME labels.
// It is adopted from labelSetToFingerprint, but avoids type conversions and memory allocations.
func FingerprintLabels(labels []model.Label) uint64 {

	if len(labels) == 0 {
		return offset64
	}

	sum := offset64
	for _, v := range labels {
		sum = hashAdd(sum, v.Key)
		sum = hashAddByte(sum, separatorByte)
		sum = hashAdd(sum, v.Value)
		sum = hashAddByte(sum, separatorByte)
	}
	return sum
}

// Fingerprint calculates a fingerprint of SORTED BY NAME labels.
// It is adopted from labelSetToFingerprint, but avoids type conversions and memory allocations.
func FingerprintLabelsCityHash(data []byte) uint64 {

	if data == nil {
		return 0
	}

	return cityhash102.CityHash64(data, uint32(len(data)))
}

// Professor Daniel J. Bernstein
func FingerprintLabelsDJBHash(data []byte) uint64 {

	if data == nil {
		return 0
	}

	var hash uint = 5381

	for i := 0; i < len(data); i++ {
		hash = (hash << 5) + hash + uint(data[i])
	}

	return uint64(hash)
}

// Javascript port
func FingerprintLabelsDJBHashPrometheus(data []byte) uint32 {

	if data == nil {
		return 0
	}

	var hash int32 = 5381

	for i := len(data) - 1; i > -1; i-- {
		hash = (hash * 33) ^ int32(uint16(data[i]))
	}
	return uint32(hash)
}
