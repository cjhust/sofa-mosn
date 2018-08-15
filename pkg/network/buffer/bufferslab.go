/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package buffer

import (
	"sync"
	"time"
	"github.com/alipay/sofa-mosn/pkg/log"
)

const minShift = 3
const maxShift = 15
const errSlot = -1

type BufferSlab struct {
	minShift int
	minSize  int
	maxSize  int

	pool     []*bufferSlot
}

type bufferSlot struct {
	defaultSize int
	pool        sync.Pool
}

func NewBufferSlab() *BufferSlab {
	p := &BufferSlab{
		minShift: minShift,
		minSize:  1 << minShift,
		maxSize:  1 << maxShift,
	}
	for i := 0; i <= maxShift-minShift; i++ {
		slab := &bufferSlot{
			defaultSize: 1 << (uint)(i+minShift),
		}
		p.pool = append(p.pool, slab)
	}

	return p
}

func (p *BufferSlab) slot(size int) int {
	if size > p.maxSize {
		return errSlot
	}
	slot := 0
	shift := 0
	if size > p.minSize {
		size--
		for size > 0 {
			size = size >> 1
			shift++
		}
		slot = shift - p.minShift
	}

	return slot
}

func newBytes(size int) []byte {
	return make([]byte, size)
}

var newb, reuse, give int
var now = time.Now().Add(time.Second)
func (p *BufferSlab) Take(size int) []byte {
	if now.Before(time.Now()) {
		now = time.Now().Add(time.Second)
		log.DefaultLogger.Errorf("newb = %d, reuse = %d, give = %d", newb, reuse, give)
	}
	slot := p.slot(size)
	if slot == errSlot {
		newb++
		return newBytes(size)
	}
	v := p.pool[slot].pool.Get()
	if v == nil {
		newb++
		return newBytes(p.pool[slot].defaultSize)
	}
	reuse++
	return v.([]byte)
}

func (p *BufferSlab) Clone(old []byte) []byte {
	size := cap(old)
	buf := p.Take(size)
	copy(buf, old)
	return buf
}

func (p *BufferSlab) Copy(old []byte) []byte {
	size := cap(old) * 2
	buf := p.Take(size)
	copy(buf, old)
	p.Give(old)
	return buf
}

func (p *BufferSlab) Give(buf []byte) {
	size := cap(buf)
	slot := p.slot(size)
	if slot == errSlot {
		return
	}

	if size != int(p.pool[slot].defaultSize) {
		return
	}

	give++
	p.pool[slot].pool.Put(buf[:size])
}
