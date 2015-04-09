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

package gl

import (
	"unsafe"

	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

const (
	ArrayBuffer BufferTarget = gl.ARRAY_BUFFER

	StaticDraw  BufferUsage = gl.STATIC_DRAW
	DynamicDraw BufferUsage = gl.DYNAMIC_DRAW
)

type Buffer struct {
	internal uint32
}

type BufferTarget uint32
type BufferUsage uint32

func CreateBuffer() Buffer {
	var buffer Buffer
	gl.GenBuffers(1, &buffer.internal)
	return buffer
}

var (
	currentBuffer       Buffer
	currentBufferTarget BufferTarget
)

func (b Buffer) Bind(target BufferTarget) {
	if currentBuffer == b && currentBufferTarget == target {
		return
	}
	gl.BindBuffer(uint32(target), b.internal)
	currentBuffer = b
	currentBufferTarget = target
}

func (b Buffer) Data(data []byte, usage BufferUsage) {
	if currentBuffer != b {
		panic("buffer not bound")
	}
	var ptr unsafe.Pointer
	if len(data) != 0 {
		ptr = gl.Ptr(data)
	}
	gl.BufferData(uint32(currentBufferTarget), len(data), ptr, uint32(usage))
}

func (b Buffer) SubData(offset int, data []byte) {
	if currentBuffer != b {
		panic("buffer not bound")
	}
	var ptr unsafe.Pointer
	if len(data) != 0 {
		ptr = gl.Ptr(data)
	}
	gl.BufferSubData(uint32(currentBufferTarget), offset, len(data), ptr)
}

func (b *Buffer) Delete() {
	gl.DeleteBuffers(1, &b.internal)
	if currentBuffer == *b {
		currentBuffer = Buffer{}
	}
	b.internal = 0
}

func (b Buffer) IsValid() bool {
	return b.internal != 0
}
