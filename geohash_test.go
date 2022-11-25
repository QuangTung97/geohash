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

func TestGeohash_Left_And_Right(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: -17.3218,
		Lon: -45.0434,
	}, 5)
	assert.Equal(t, "6uzvr", h.String())

	assert.Equal(t, "6uzvq", h.Left().String())
	assert.Equal(t, "7hbj2", h.Right().String())

	h = ComputeGeohash(Pos{
		Lat: -17.3218,
		Lon: -45.0200,
	}, 5)
	assert.Equal(t, "6uzvr", h.String())

	h = ComputeGeohash(Pos{
		Lat: -89.97802734,
		Lon: -179.97802734,
	}, 5)
	assert.Equal(t, "00000", h.String())
	assert.Equal(t, "pbpbp", h.Left().String())
	assert.Equal(t, "00001", h.Right().String())
}

func TestGeohash_Left_And_Right_Case_2(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: 89.97802734,
		Lon: 179.97802734,
	}, 5)
	assert.Equal(t, "zzzzz", h.String())
	assert.Equal(t, "zzzzy", h.Left().String())
	assert.Equal(t, "bpbpb", h.Right().String())
}

func TestGeohash_Top_And_Bot(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: -17.3218,
		Lon: -45.0434,
	}, 5)
	assert.Equal(t, "6uzvr", h.String())

	assert.Equal(t, "6uzvx", h.Top().String())
	assert.Equal(t, "6uzvp", h.Bottom().String())

	h = ComputeGeohash(Pos{
		Lat: -17.3218,
		Lon: -45.0200,
	}, 5)
	assert.Equal(t, "6uzvr", h.String())

	h = ComputeGeohash(Pos{
		Lat: -89.97802734,
		Lon: -179.97802734,
	}, 5)
	assert.Equal(t, "00002", h.Top().String())
	assert.Equal(t, "bpbpb", h.Bottom().String())
}

func TestGeohash_Top_And_Bottom_Case_2(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: 89.97802734,
		Lon: 179.97802734,
	}, 5)
	assert.Equal(t, "zzzzz", h.String())
	assert.Equal(t, "pbpbp", h.Top().String())
	assert.Equal(t, "zzzzx", h.Bottom().String())
}

func TestGeohash_Rec(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: 0,
		Lon: 0,
	}, 5)
	assert.Equal(t, "s0000", h.String())
	assert.Equal(t, Rectangle{
		BottomLeft: Pos{
			Lat: 0,
			Lon: 0,
		},
		BottomRight: Pos{
			Lat: 0,
			Lon: 0.0439453125,
		},
		TopLeft: Pos{
			Lat: 0.0439453125,
			Lon: 0,
		},
		TopRight: Pos{
			Lat: 0.0439453125,
			Lon: 0.0439453125,
		},
	}, h.Rec())
}

func TestNearbyNext(t *testing.T) {
	offset := posOffset{
		lat: 0,
		lon: 1,
	}
	var ok bool

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: 1,
		lon: 1,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: 1,
		lon: 0,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: 1,
		lon: -1,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: 0,
		lon: -1,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: -1,
		lon: -1,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: -1,
		lon: 0,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, true, ok)
	assert.Equal(t, posOffset{
		lat: -1,
		lon: 1,
	}, offset)

	offset, ok = nearbyNext(offset, 1)
	assert.Equal(t, false, ok)
	assert.Equal(t, posOffset{}, offset)
}

func TestNearbyNext__Radius_2(t *testing.T) {
	var offsets []posOffset

	offset := posOffset{
		lat: 0,
		lon: 2,
	}
	offsets = append(offsets, offset)

	for {
		var ok bool
		offset, ok = nearbyNext(offset, 2)
		if !ok {
			break
		}
		offsets = append(offsets, offset)
	}

	assert.Equal(t, []posOffset{
		{lat: 0, lon: 2},
		{lat: 1, lon: 2},
		{lat: 2, lon: 2},
		{lat: 2, lon: 1},
		{lat: 2, lon: 0},
		{lat: 2, lon: -1},
		{lat: 2, lon: -2},
		{lat: 1, lon: -2},
		{lat: 0, lon: -2},
		{lat: -1, lon: -2},
		{lat: -2, lon: -2},
		{lat: -2, lon: -1},
		{lat: -2, lon: 0},
		{lat: -2, lon: 1},
		{lat: -2, lon: 2},
		{lat: -1, lon: 2},
	}, offsets)
}

func TestNearbyGeohashs(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: 0.7,
		Lon: 0.7,
	}, 3)
	assert.Equal(t, "s00", h.String())

	origin := Pos{
		Lat: 0.7,
		Lon: 0.7,
	}

	rec := h.Rec()
	assert.Equal(t, 110.56042392519969, haversineDistance(origin, rec.TopLeft))

	hashList := NearbyGeohashs(origin, 20, 3)
	assert.Equal(t, []Hash{h}, hashList)

	hashList = NearbyGeohashs(Pos{
		Lat: 0.7,
		Lon: 0.7,
	}, 120, 3)
	assert.Equal(t, []Hash{
		h,
		h.Right(), h.Right().Top(),
		h.Top(), h.Top().Left(),
		h.Left(), h.Bottom().Left(),
		h.Bottom(), h.Bottom().Right(),
	}, hashList)

	hashList = NearbyGeohashs(Pos{
		Lat: 0.7,
		Lon: 0.7,
	}, 80, 3)
	assert.Equal(t, []Hash{
		h,
		h.Right(),
		h.Top(),
		h.Left(),
		h.Bottom(),
	}, hashList)
}

func TestNearestTopEdge(t *testing.T) {
	h := ComputeGeohash(Pos{
		Lat: 0.7,
		Lon: 0.7,
	}, 3)
	assert.Equal(t, "s00", h.String())

	const size = 1.40625

	assert.Equal(t, Rectangle{
		TopRight: Pos{
			Lat: size,
			Lon: size,
		},
		TopLeft: Pos{
			Lat: size,
			Lon: 0,
		},
		BottomRight: Pos{
			Lat: 0,
			Lon: size,
		},
	}, h.Rec())

	t.Run("inside-range", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: 0.3,
		}, nearestTopEdge(Pos{
			Lat: 10,
			Lon: 0.3,
		}, h.Rec()))
	})

	t.Run("outside-left", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: 0,
		}, nearestTopEdge(Pos{
			Lat: 10,
			Lon: -0.3,
		}, h.Rec()))
	})

	t.Run("outside-right", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: size,
		}, nearestTopEdge(Pos{
			Lat: 10,
			Lon: size + 3.0,
		}, h.Rec()))
	})
}

func TestNearestBottomEdge(t *testing.T) {
	const size = 1.40625

	h := ComputeGeohash(Pos{
		Lat: size + 1,
		Lon: size + 1,
	}, 3)
	assert.Equal(t, "s03", h.String())

	assert.Equal(t, Rectangle{
		BottomLeft: Pos{
			Lat: size,
			Lon: size,
		},
		TopRight: Pos{
			Lat: size + size,
			Lon: size + size,
		},
		TopLeft: Pos{
			Lat: size + size,
			Lon: size,
		},
		BottomRight: Pos{
			Lat: size,
			Lon: size + size,
		},
	}, h.Rec())

	t.Run("inside-range", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: size + 0.3,
		}, nearestBottomEdge(Pos{
			Lat: 10,
			Lon: size + 0.3,
		}, h.Rec()))
	})

	t.Run("outside-left", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: size,
		}, nearestBottomEdge(Pos{
			Lat: 10,
			Lon: size - 0.3,
		}, h.Rec()))
	})

	t.Run("outside-right", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: size + size,
		}, nearestBottomEdge(Pos{
			Lat: 10,
			Lon: size + size + 3.0,
		}, h.Rec()))
	})
}

func TestNearestLeftEdge(t *testing.T) {
	const size = 1.40625

	h := ComputeGeohash(Pos{
		Lat: size + 1,
		Lon: size + 1,
	}, 3)
	assert.Equal(t, "s03", h.String())

	assert.Equal(t, Rectangle{
		BottomLeft: Pos{
			Lat: size,
			Lon: size,
		},
		TopRight: Pos{
			Lat: size + size,
			Lon: size + size,
		},
		TopLeft: Pos{
			Lat: size + size,
			Lon: size,
		},
		BottomRight: Pos{
			Lat: size,
			Lon: size + size,
		},
	}, h.Rec())

	t.Run("inside-range", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: 1.7256124742605874,
			Lon: size,
		}, nearestLeftEdge(Pos{
			Lat: size + 0.3,
			Lon: 10,
		}, h.Rec()))
	})

	t.Run("outside-bottom", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: size,
		}, nearestLeftEdge(Pos{
			Lat: size - 0.3,
			Lon: 10,
		}, h.Rec()))
	})

	t.Run("outside-right", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size + size,
			Lon: size,
		}, nearestLeftEdge(Pos{
			Lat: size + size + 3.0,
			Lon: 10,
		}, h.Rec()))
	})
}

func TestNearestRightEdge(t *testing.T) {
	const size = 1.40625

	h := ComputeGeohash(Pos{
		Lat: size + 1,
		Lon: size + 1,
	}, 3)
	assert.Equal(t, "s03", h.String())

	assert.Equal(t, Rectangle{
		BottomLeft: Pos{
			Lat: size,
			Lon: size,
		},
		TopRight: Pos{
			Lat: size + size,
			Lon: size + size,
		},
		TopLeft: Pos{
			Lat: size + size,
			Lon: size,
		},
		BottomRight: Pos{
			Lat: size,
			Lon: size + size,
		},
	}, h.Rec())

	t.Run("inside-range", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: 1.7197557847302827,
			Lon: size + size,
		}, nearestRightEdge(Pos{
			Lat: size + 0.3,
			Lon: 10,
		}, h.Rec()))
	})

	t.Run("outside-bottom", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size,
			Lon: size + size,
		}, nearestRightEdge(Pos{
			Lat: size - 0.3,
			Lon: 10,
		}, h.Rec()))
	})

	t.Run("outside-right", func(t *testing.T) {
		assert.Equal(t, Pos{
			Lat: size + size,
			Lon: size + size,
		}, nearestRightEdge(Pos{
			Lat: size + size + 3.0,
			Lon: 10,
		}, h.Rec()))
	})
}
