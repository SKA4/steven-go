// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package render

import (
	"math"

	"github.com/thinkofdeath/steven/native"
	"github.com/thinkofdeath/steven/render/gl"
)

const (
	uiWidth, uiHeight = 800, 480
)

var (
	uiState = struct {
		program      gl.Program
		shader       *uiShader
		array        gl.VertexArray
		buffer       gl.Buffer
		count        int
		data         []byte
		prevSize     int
		elements     []*UIElement
		freeElements []*UIElement
		lastID       int
		elementCount int
	}{
		prevSize: -1,
	}
)

func initUI() {
	uiState.program = CreateProgram(vertexUI, fragmentUI)
	uiState.shader = &uiShader{}
	InitStruct(uiState.shader, uiState.program)

	uiState.array = gl.CreateVertexArray()
	uiState.array.Bind()
	uiState.buffer = gl.CreateBuffer()
	uiState.buffer.Bind(gl.ArrayBuffer)
	uiState.shader.Position.Enable()
	uiState.shader.TextureInfo.Enable()
	uiState.shader.TextureOffset.Enable()
	uiState.shader.Color.Enable()
	uiState.shader.Position.Pointer(3, gl.Float, false, 28, 0)
	uiState.shader.TextureInfo.Pointer(4, gl.UnsignedShort, false, 28, 12)
	uiState.shader.TextureOffset.Pointer(2, gl.Short, false, 28, 20)
	uiState.shader.Color.Pointer(4, gl.UnsignedByte, true, 28, 24)

	AddUIElement(GetTexture("stone"), 400-2, 240-2, 4, 4, 0, 0, 16, 16)
}

func drawUI() {
	// Redraw everything
	uiState.count = 0
	uiState.data = uiState.data[:0]
	for _, e := range uiState.elements {
		e.draw()
	}

	// Prevent clipping with the world
	gl.Clear(gl.DepthBufferBit)
	gl.Enable(gl.Blend)

	uiState.program.Use()
	uiState.shader.Texture.Int(0)
	if uiState.count > 0 {
		uiState.array.Bind()
		uiState.buffer.Bind(gl.ArrayBuffer)
		if len(uiState.data) > uiState.prevSize {
			uiState.prevSize = len(uiState.data)
			uiState.buffer.Data(uiState.data, gl.DynamicDraw)
		} else {
			target := uiState.buffer.Map(gl.WriteOnly, len(uiState.data))
			copy(target, uiState.data)
			uiState.buffer.Unmap()
		}
		gl.DrawArrays(gl.Triangles, 0, uiState.count)
	}
	gl.Disable(gl.Blend)
}

// UIElement is a single element on the screen. It is a rectangle
// with a texture and a tint.
type UIElement struct {
	free bool

	X, Y, W, H         float64
	DepthIndex         float64
	TX, TY, TW, TH     uint16
	TOffsetX, TOffsetY int16
	TSizeW, TSizeH     int16
	R, G, B, A         byte
}

// AddUIElement creates and adds a single ui element onto the screen.
func AddUIElement(tex *TextureInfo, x, y, width, height float64, tx, ty, tw, th int) *UIElement {
	var e *UIElement
	if len(uiState.freeElements) == 0 {
		e = &UIElement{}
		uiState.elements = append(uiState.elements, e)
	} else {
		l := len(uiState.freeElements)
		e = uiState.freeElements[l-1]
		uiState.freeElements = uiState.freeElements[:l-1]
	}
	// (Re)set the information for the element
	e.X = x / uiWidth
	e.Y = y / uiHeight
	e.W = width / uiWidth
	e.H = height / uiHeight
	e.TX = uint16(tex.X)
	e.TY = uint16(tex.Y + tex.Atlas*AtlasSize)
	e.TW = uint16(tex.Width)
	e.TH = uint16(tex.Height)
	e.TOffsetX = int16(tx * 16)
	e.TOffsetY = int16(ty * 16)
	e.TSizeW = int16(tw * 16)
	e.TSizeH = int16(th * 16)
	e.R = 255
	e.G = 255
	e.B = 255
	e.A = 255
	e.DepthIndex = -float64(uiState.elementCount) / float64(math.MaxInt16)
	uiState.elementCount++
	e.free = false
	return e
}

// Shift moves the element by the passed amounts.
func (u *UIElement) Shift(x, y float64) {
	u.X += x / uiWidth
	u.Y += y / uiHeight
}

// Alpha changes the alpha of this element
func (u *UIElement) Alpha(a float64) {
	if a > 1.0 {
		a = 1.0
	}
	u.A = byte(255.0 * a)
}

// Free removes the element from the screen. This may be reused
// so this element should be considered invalid after this call.
func (u *UIElement) Free() {
	if u.free {
		return
	}
	u.free = true
	uiState.freeElements = append(uiState.freeElements, u)
	uiState.elementCount--
}

func (u *UIElement) draw() {
	if u.free {
		return
	}
	u.appendVertex(u.X, u.Y, u.TOffsetX, u.TOffsetY)
	u.appendVertex(u.X+u.W, u.Y, u.TOffsetX+u.TSizeW, u.TOffsetY)
	u.appendVertex(u.X, u.Y+u.H, u.TOffsetX, u.TOffsetY+u.TSizeH)

	u.appendVertex(u.X+u.W, u.Y+u.H, u.TOffsetX+u.TSizeW, u.TOffsetY+u.TSizeH)
	u.appendVertex(u.X, u.Y+u.H, u.TOffsetX, u.TOffsetY+u.TSizeH)
	u.appendVertex(u.X+u.W, u.Y, u.TOffsetX+u.TSizeW, u.TOffsetY)
}

func (u *UIElement) appendVertex(x, y float64, tx, ty int16) {
	uiState.count++
	uiState.data = appendFloat(uiState.data, float32(x))
	uiState.data = appendFloat(uiState.data, float32(y))
	uiState.data = appendFloat(uiState.data, float32(u.DepthIndex))
	uiState.data = appendUnsignedShort(uiState.data, u.TX)
	uiState.data = appendUnsignedShort(uiState.data, u.TY)
	uiState.data = appendUnsignedShort(uiState.data, u.TW)
	uiState.data = appendUnsignedShort(uiState.data, u.TH)
	uiState.data = appendShort(uiState.data, tx)
	uiState.data = appendShort(uiState.data, ty)
	uiState.data = appendUnsignedByte(uiState.data, u.R)
	uiState.data = appendUnsignedByte(uiState.data, u.G)
	uiState.data = appendUnsignedByte(uiState.data, u.B)
	uiState.data = appendUnsignedByte(uiState.data, u.A)
}

func appendUnsignedByte(data []byte, i byte) []byte {
	return append(data, i)
}

func appendByte(data []byte, i int8) []byte {
	return appendUnsignedByte(data, byte(i))
}

var scratch [8]byte

func appendUnsignedShort(data []byte, i uint16) []byte {
	d := scratch[:2]
	native.Order.PutUint16(d, i)
	return append(data, d...)
}

func appendShort(data []byte, i int16) []byte {
	return appendUnsignedShort(data, uint16(i))
}

func appendFloat(data []byte, f float32) []byte {
	d := scratch[:4]
	i := math.Float32bits(f)
	native.Order.PutUint32(d, i)
	return append(data, d...)
}
