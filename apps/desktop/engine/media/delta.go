package media

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

const (
	deltaMagic      = "NDF1"
	deltaVersion    = 1
	defaultTileSize = 32
	modeCompressed  = 0
	modeSolid       = 1
	modeXOR         = 2
)

// DeltaEncoder turns successive RGBA frames into keyframes and dirty tile updates.
type DeltaEncoder struct {
	TileSize int
	previous *image.RGBA
}

// NewDeltaEncoder creates an encoder using 32 by 32 tiles unless tileSize is positive.
func NewDeltaEncoder(tileSize int) *DeltaEncoder {
	if tileSize <= 0 {
		tileSize = defaultTileSize
	}
	return &DeltaEncoder{TileSize: tileSize}
}

// Encode returns one complete delta packet. The first frame and resized frames are keyframes.
func (e *DeltaEncoder) Encode(frame *image.RGBA) ([]byte, error) {
	if frame == nil || frame.Bounds().Empty() {
		return nil, fmt.Errorf("media: cannot encode an empty frame")
	}
	bounds := frame.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	keyframe := e.previous == nil || e.previous.Bounds() != bounds
	if keyframe {
		e.previous = image.NewRGBA(bounds)
	}

	tiles := make([]encodedTile, 0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += e.TileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += e.TileSize {
			tile := tileRect(bounds, x, y, e.TileSize)
			current := rgbaBytes(frame, tile)
			before := rgbaBytes(e.previous, tile)
			if !keyframe && bytes.Equal(current, before) {
				continue
			}

			mode := byte(modeCompressed)
			payload := current
			if color, ok := solidColor(current); ok {
				mode = modeSolid
				payload = color[:]
			} else if !keyframe {
				mode = modeXOR
				payload = xorBytes(current, before)
			}
			compressed, err := compress(payload)
			if err != nil {
				return nil, err
			}
			tiles = append(tiles, encodedTile{x: x - bounds.Min.X, y: y - bounds.Min.Y, width: tile.Dx(), height: tile.Dy(), mode: mode, payload: compressed})
		}
	}

	var out bytes.Buffer
	out.WriteString(deltaMagic)
	out.WriteByte(deltaVersion)
	if keyframe {
		out.WriteByte(1)
	} else {
		out.WriteByte(0)
	}
	writeUint32(&out, uint32(width))
	writeUint32(&out, uint32(height))
	writeUint16(&out, uint16(e.TileSize))
	writeUint32(&out, uint32(len(tiles)))
	for _, tile := range tiles {
		writeUint16(&out, uint16(tile.x))
		writeUint16(&out, uint16(tile.y))
		writeUint16(&out, uint16(tile.width))
		writeUint16(&out, uint16(tile.height))
		out.WriteByte(tile.mode)
		writeUint32(&out, uint32(len(tile.payload)))
		out.Write(tile.payload)
	}

	copy(e.previous.Pix, frame.Pix)
	return out.Bytes(), nil
}

type encodedTile struct {
	x, y, width, height int
	mode                byte
	payload             []byte
}

func tileRect(bounds image.Rectangle, x, y, size int) image.Rectangle {
	maxX, maxY := x+size, y+size
	if maxX > bounds.Max.X {
		maxX = bounds.Max.X
	}
	if maxY > bounds.Max.Y {
		maxY = bounds.Max.Y
	}
	return image.Rect(x, y, maxX, maxY)
}

func rgbaBytes(img *image.RGBA, rect image.Rectangle) []byte {
	result := make([]byte, 0, rect.Dx()*rect.Dy()*4)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		start := img.PixOffset(rect.Min.X, y)
		result = append(result, img.Pix[start:start+rect.Dx()*4]...)
	}
	return result
}

func solidColor(data []byte) ([4]byte, bool) {
	var color [4]byte
	if len(data) < 4 {
		return color, false
	}
	copy(color[:], data[:4])
	for i := 4; i < len(data); i += 4 {
		if !bytes.Equal(data[i:i+4], color[:]) {
			return color, false
		}
	}
	return color, true
}

func xorBytes(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func compress(data []byte) ([]byte, error) {
	var out bytes.Buffer
	writer := zlib.NewWriter(&out)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func writeUint16(out *bytes.Buffer, value uint16) { _ = binary.Write(out, binary.BigEndian, value) }
func writeUint32(out *bytes.Buffer, value uint32) { _ = binary.Write(out, binary.BigEndian, value) }
