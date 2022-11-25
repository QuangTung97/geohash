package geohash

import (
	"encoding/binary"
	"github.com/QuangTung97/haversine"
	"math"
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

// Rectangle represents all corners
type Rectangle struct {
	BottomLeft  Pos
	BottomRight Pos
	TopLeft     Pos
	TopRight    Pos
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

func (h Hash) Left() Hash {
	return h.addOffset(posOffset{
		lat: 0,
		lon: -1,
	})
}

func (h Hash) Right() Hash {
	return h.addOffset(posOffset{
		lat: 0,
		lon: 1,
	})
}

func (h Hash) Top() Hash {
	return h.addOffset(posOffset{
		lat: 1,
		lon: 0,
	})
}

func (h Hash) Bottom() Hash {
	return h.addOffset(posOffset{
		lat: -1,
		lon: 0,
	})
}

// Pos returns the bottom left position of this geohash
func (h Hash) Pos() Pos {
	bitCount := h.precision * 5
	latPrecision := bitCount >> 1
	lonPrecision := bitCount - latPrecision

	lat := bitsToLat(h.lat, uint32(1<<latPrecision))
	lon := bitsToLon(h.lon, uint32(1<<lonPrecision))

	return Pos{Lat: lat, Lon: lon}
}

// Rec returns 4 corners of this geohash
func (h Hash) Rec() Rectangle {
	top := h.Top()
	return Rectangle{
		BottomLeft:  h.Pos(),
		BottomRight: h.Right().Pos(),
		TopLeft:     top.Pos(),
		TopRight:    top.Right().Pos(),
	}
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

func bitsToLon(bits uint32, multiplier uint32) float64 {
	return float64(bits)*360/float64(multiplier) - 180
}

func latToBits(lat float64, multiplier uint32) uint32 {
	return uint32((lat + 90) * float64(multiplier) / 180)
}

func bitsToLat(bits uint32, multiplier uint32) float64 {
	return float64(bits)*180/float64(multiplier) - 90
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

func (p Pos) toHaversine() haversine.Pos {
	return haversine.Pos{
		Lat: p.Lat,
		Lon: p.Lon,
	}
}

func haversineDistance(a, b Pos) float64 {
	return haversine.DistanceEarth(a.toHaversine(), b.toHaversine())
}

func minDistanceToGeohash(origin Pos, hash Hash) float64 {
	rec := hash.Rec()

	minDistance := haversineDistance(origin, nearestLeftEdge(origin, rec))

	d := haversineDistance(origin, nearestRightEdge(origin, rec))
	minDistance = math.Min(minDistance, d)

	d = haversineDistance(origin, nearestTopEdge(origin, rec))
	minDistance = math.Min(minDistance, d)

	d = haversineDistance(origin, nearestBottomEdge(origin, rec))
	minDistance = math.Min(minDistance, d)

	return minDistance
}

func (h Hash) addOffset(offset posOffset) Hash {
	bitCount := h.precision * 5
	latPrecision := bitCount >> 1
	lonPrecision := bitCount - latPrecision

	latMask := uint32((1 << latPrecision) - 1)
	lonMask := uint32((1 << lonPrecision) - 1)

	h.lat = (h.lat + uint32(int(latMask)+1+offset.lat)) & latMask
	h.lon = (h.lon + uint32(int(lonMask)+1+offset.lon)) & lonMask

	return h
}

// NearbyGeohashs computes nearby geohashs, radius is in km
func NearbyGeohashs(origin Pos, radius float64, precision uint32) []Hash {
	h := ComputeGeohash(origin, precision)

	result := []Hash{h}

	for distance := 1; ; distance++ {
		continuing := false

		offset := posOffset{lat: 0, lon: distance}
		ok := true
		for ; ok; offset, ok = nearbyNext(offset, distance) {
			newHash := h.addOffset(offset)

			d := minDistanceToGeohash(origin, newHash)
			if d > radius {
				continue
			}

			continuing = true
			result = append(result, newHash)
		}

		if !continuing {
			return result
		}
	}
}

type posOffset struct {
	lat int
	lon int
}

func (a posOffset) add(b posOffset) posOffset {
	return posOffset{
		lat: a.lat + b.lat,
		lon: a.lon + b.lon,
	}
}

func directionOfOffset(off posOffset, radius int) posOffset {
	if off.lon == radius {
		return posOffset{lat: 1, lon: 0}
	}
	if off.lat == radius {
		return posOffset{lat: 0, lon: -1}
	}
	if off.lon == -radius {
		return posOffset{lat: -1, lon: 0}
	}
	return posOffset{lat: 0, lon: 1}
}

func rotateDirection(off posOffset) posOffset {
	return posOffset{lat: off.lon, lon: -off.lat}
}

func nearbyNext(offset posOffset, radius int) (posOffset, bool) {
	direction := directionOfOffset(offset, radius)

	if offset.lat == radius && offset.lon == radius {
		direction = rotateDirection(direction)
	}

	if offset.lat == radius && offset.lon == -radius {
		direction = rotateDirection(direction)
	}

	if offset.lat == -radius && offset.lon == -radius {
		direction = rotateDirection(direction)
	}

	offset = offset.add(direction)
	if offset.lat == 0 && offset.lon == radius {
		return posOffset{}, false
	}

	return offset, true
}

func nearestHorizontalEdge(pos Pos, lat float64, rec Rectangle) Pos {
	lon := pos.Lon
	if lon < rec.TopLeft.Lon {
		lon = rec.TopLeft.Lon
	} else if lon > rec.TopRight.Lon {
		lon = rec.TopRight.Lon
	}

	return Pos{
		Lat: lat,
		Lon: lon,
	}
}

func nearestTopEdge(pos Pos, rec Rectangle) Pos {
	return nearestHorizontalEdge(pos, rec.TopRight.Lat, rec)
}

func nearestBottomEdge(pos Pos, rec Rectangle) Pos {
	return nearestHorizontalEdge(pos, rec.BottomRight.Lat, rec)
}

func nearestVerticalEdge(pos Pos, lon float64, rec Rectangle) Pos {
	lat := haversine.MinLatDistance(pos.toHaversine(), lon)

	minLat := rec.BottomLeft.Lat
	maxLat := rec.TopLeft.Lat
	if lat < minLat {
		lat = minLat
	} else if lat > maxLat {
		lat = maxLat
	}

	return Pos{
		Lat: lat,
		Lon: lon,
	}
}

func nearestLeftEdge(pos Pos, rec Rectangle) Pos {
	return nearestVerticalEdge(pos, rec.TopLeft.Lon, rec)
}

func nearestRightEdge(pos Pos, rec Rectangle) Pos {
	return nearestVerticalEdge(pos, rec.TopRight.Lon, rec)
}
