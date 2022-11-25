package geohash

import (
	"encoding/binary"
)

type Pos struct {
	Lat float64 // in degree
	Lon float64 // in degree
}

// Hash ...
type Hash struct {
	precision uint32
	lat       uint32 // lat bits
	lon       uint32 // lon bits
}

// spacingByte only max 4 bits
func spacingByte(a uint8) uint8 {
	result := a & 0b1
	result |= (a & 0b10) << 1
	result |= (a & 0b100) << 2
	result |= (a & 0b1000) << 3
	return result
}

func spacing(bits uint64, count uint32) uint64 {
	bytesCount := (count + 3) / 4

	bytes := [8]uint8{}
	for index := uint32(0); index < bytesCount; index++ {
		val := uint8(bits & 0xff)
		bits >>= 4

		newVal := spacingByte(val)
		bytes[index] = newVal
	}

	return binary.LittleEndian.Uint64(bytes[:])
}

func (h Hash) String() string {
	bitCount := h.precision * 5
	latPrecision := bitCount >> 1
	lonPrecision := bitCount - latPrecision

	latBits := spacing(uint64(h.lat), latPrecision)
	lonBits := spacing(uint64(h.lon), lonPrecision)

	var resultHash uint64

	if latPrecision == lonPrecision {
		resultHash = latBits | (lonBits << 1)
	} else {
		resultHash = (latBits << 1) | lonBits
	}

	var resultBytes [12]uint8
	for index := uint32(0); index < h.precision; index++ {
		resultBytes[index] = uint8(resultHash & 0b11111)
		resultHash >>= 5
	}

	stringBytes := make([]byte, 0, 12)
	for i := int(h.precision - 1); i >= 0; i-- {
		b := encoding[resultBytes[i]]
		stringBytes = append(stringBytes, b)
	}

	return string(stringBytes)
}

var encoding = []byte{
	'0', '1', '2', '3',
	'4', '5', '6', '7',
	'8', '9', 'b', 'c',
	'd', 'e', 'f', 'g',
	'h', 'j', 'k', 'm',
	'n', 'p', 'q', 'r',
	's', 't', 'u', 'v',
	'w', 'x', 'y', 'z',
}

func lonToBits(lon float64, multiplier uint32) uint32 {
	return uint32((lon + 180) * float64(multiplier) / 360)
}

func latToBits(lat float64, multiplier uint32) uint32 {
	return uint32((lat + 90) * float64(multiplier) / 180)
}

// ComputeGeohash support precision <= 12
func ComputeGeohash(pos Pos, precision uint32) Hash {
	bitCount := precision * 5
	latPrecision := bitCount >> 1
	lonPrecision := bitCount - latPrecision

	lat := latToBits(pos.Lat, uint32(1<<latPrecision))
	lon := lonToBits(pos.Lon, uint32(1<<lonPrecision))

	return Hash{
		precision: precision,
		lat:       lat,
		lon:       lon,
	}
}
