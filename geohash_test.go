package geohash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpacing(t *testing.T) {
	result := spacing(0b11, 2)
	assert.Equal(t, uint64(0b101), result)

	result = spacing(0b101, 3)
	assert.Equal(t, uint64(0b10001), result)

	result = spacing(0b1111111, 7)
	assert.Equal(t, uint64(0b1010101010101), result)

	result = spacing(0b1111001, 7)
	assert.Equal(t, uint64(0b1010101000001), result)

	result = spacing(0b11111001, 8)
	assert.Equal(t, uint64(0b101010101000001), result)
}

func TestComputeGeohash(t *testing.T) {
	result := ComputeGeohash(Pos{
		Lat: 48.669,
		Lon: 22.445,
	}, 3).String()
	assert.Equal(t, "u2x", result)

	result = ComputeGeohash(Pos{
		Lat: 48.669,
		Lon: 22.445,
	}, 4).String()
	assert.Equal(t, "u2xu", result)

	result = ComputeGeohash(Pos{
		Lat: 48.669,
		Lon: 22.445,
	}, 5).String()
	assert.Equal(t, "u2xuy", result)

	result = ComputeGeohash(Pos{
		Lat: 48.669,
		Lon: 22.445,
	}, 6).String()
	assert.Equal(t, "u2xuye", result)

	result = ComputeGeohash(Pos{
		Lat: 48.66746,
		Lon: 22.44043,
	}, 7).String()
	assert.Equal(t, "u2xuyes", result)

	result = ComputeGeohash(Pos{
		Lat: 48.66746,
		Lon: 22.44043,
	}, 8).String()
	assert.Equal(t, "u2xuyess", result)
}

func TestComputeGeohash_Case2(t *testing.T) {
	result := ComputeGeohash(Pos{
		Lat: -10.669,
		Lon: 12.445,
	}, 3).String()
	assert.Equal(t, "kq0", result)

	result = ComputeGeohash(Pos{
		Lat: -10.669,
		Lon: 12.445,
	}, 4).String()
	assert.Equal(t, "kq0g", result)

	result = ComputeGeohash(Pos{
		Lat: -10.669,
		Lon: 12.445,
	}, 5).String()
	assert.Equal(t, "kq0g7", result)

	result = ComputeGeohash(Pos{
		Lat: -10.669,
		Lon: 12.445,
	}, 6).String()
	assert.Equal(t, "kq0g71", result)

	result = ComputeGeohash(Pos{
		Lat: -10.6698,
		Lon: 12.4457,
	}, 7).String()
	assert.Equal(t, "kq0g71w", result)
}

func TestComputeGeohash_Case3(t *testing.T) {
	result := ComputeGeohash(Pos{
		Lat: 28.3218,
		Lon: -62.0434,
	}, 7).String()
	assert.Equal(t, "dt5ch5v", result)
}

func TestComputeGeohash_Case4(t *testing.T) {
	result := ComputeGeohash(Pos{
		Lat: -17.3218,
		Lon: -45.0434,
	}, 7).String()
	assert.Equal(t, "6uzvrn8", result)
}
