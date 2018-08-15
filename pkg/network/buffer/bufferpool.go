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
	"github.com/alipay/sofa-mosn/pkg/types"
	"sync"
	"context"
	"runtime"
	"sync/atomic"
)

const poolSizeMax = 64

var (
	index uint32
	poolSize = runtime.NumCPU()
    bufferPools [poolSizeMax]bufferPool
)

type bufferPool struct {
	sync.Map
}

type bufferDefaultMode struct {
	ctx types.BufferCommonCtx
	sync.Pool
}

func (mode *bufferDefaultMode) Take(i interface{}) (value interface{}) {
	value = mode.Get()
	if value == nil {
		value = mode.ctx.New(i)
	}
	return
}

func (mode *bufferDefaultMode) Give(value interface{}) {
	mode.ctx.Reset(value)
	mode.Put(value)
}

type BufferCtx struct {
	bufferPool *bufferPool
	clean      []*bufferClean
}

type bufferClean struct {
	ctx   types.BufferCommonCtx
	value interface{}
}

func NewBufferPoolContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, types.ContextKeyBufferPoolCtx, NewBufferCtx())
}

func NewBufferPoolRequestContext(ctx context.Context) context.Context {
	bufferCtx := BufferPoolContext(ctx)
	newctx := BufferCtxCopy(bufferCtx)
	return context.WithValue(ctx, types.ContextKeyBufferPoolCtx, newctx)
}

func CopyBufferPoolRequestContext(rctx context.Context, pctx context.Context) {
	x := BufferPoolContext(rctx)
	y := BufferPoolContext(pctx)
	BufferCtxCopyClean(x, y)
}

func NewBufferCtx() *BufferCtx {
	i := atomic.AddUint32(&index, 1)
	i = i % uint32(poolSizeMax) % uint32(poolSize)

	return &BufferCtx{
		bufferPool: &bufferPools[int(i)],
	}
}

func(ctx *BufferCtx) Len() int {
	return len(ctx.clean)
}

func (ctx *BufferCtx) Find(commonCtx types.BufferCommonCtx, i interface{}) (value interface{}) {
	for i := 0; i < len(ctx.clean); i++ {
		if commonCtx.Name() == ctx.clean[i].ctx.Name() {
			return ctx.clean[i].value
		}
	}
	return ctx.Take(commonCtx, i)
}

func (ctx *BufferCtx) GetPool(commonCtx types.BufferCommonCtx) types.BufferPoolMode {
	load, ok := ctx.bufferPool.Load(commonCtx.Name())
	if !ok {
		init := commonCtx.Init()
		if init == nil {
			init = &bufferDefaultMode{
				ctx: commonCtx,
			}
		}
		ctx.bufferPool.Store(commonCtx.Name(), init)
		load = init
	}
	mode, _ := load.(types.BufferPoolMode)
	return mode
}

func (ctx *BufferCtx) Take(commonCtx types.BufferCommonCtx, i interface{}) (value interface{}) {
	pool := ctx.GetPool(commonCtx)
	value = pool.Take(i)
	ctx.clean = append(ctx.clean, &bufferClean{commonCtx, value})
	return
}

func (ctx *BufferCtx) Give() {
	for i := 0; i < len(ctx.clean); i++ {
		clean := ctx.clean[i]
		load, _ := ctx.bufferPool.Load(clean.ctx.Name())
		pool := load.(types.BufferPoolMode)
		pool.Give(clean.value)
	}
}

func BufferCtxCopy(ctx *BufferCtx) *BufferCtx {
	newctx := new(BufferCtx)
	newctx.bufferPool = ctx.bufferPool
	newctx.clean = append(newctx.clean, ctx.clean...)
    ctx.clean = ctx.clean[:0]
	return newctx
}

func BufferCtxCopyClean(dst *BufferCtx, src *BufferCtx) {
	dst.clean = append(dst.clean, src.clean...)
	src.clean = src.clean[:0]
}

func BufferPoolContext(context context.Context) *BufferCtx {
	if context != nil && context.Value(types.ContextKeyBufferPoolCtx) != nil {
		return context.Value(types.ContextKeyBufferPoolCtx).(*BufferCtx)
	}
	return nil
}


func BufferGetPool(commonCtx types.BufferCommonCtx) types.BufferPoolMode {
	i := index % uint32(poolSizeMax) % uint32(poolSize)
	bufferPool := bufferPools[int(i)]
	load, ok := bufferPool.Load(commonCtx.Name())
	if !ok {
		init := commonCtx.Init()
		if init == nil {
			init = &bufferDefaultMode{
				ctx: commonCtx,
			}
		}
		bufferPool.Store(commonCtx.Name(), init)
		load = init
	}
	mode, _ := load.(types.BufferPoolMode)
	return mode
}