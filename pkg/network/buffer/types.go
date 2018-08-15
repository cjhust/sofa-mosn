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

import "github.com/alipay/sofa-mosn/pkg/types"

// []byte
type BufferByteCtx struct{}

func (ctx BufferByteCtx) Name() string {
	return "BufferByteCtx"
}

func (ctx BufferByteCtx) Init() types.BufferPoolMode {
	return &bufferByteMode{
		NewBufferSlab(),
	}
}

func (ctx BufferByteCtx) New(i interface{}) interface{} {
	return nil
}

func (ctx BufferByteCtx) Reset(i interface{}) {
    return
}

type bufferByteMode struct {
	 pool *BufferSlab
}

func (mode *bufferByteMode) Take(i interface{}) interface{} {
	return mode.pool.Take(i.(int))
}

func (mode *bufferByteMode) Give(value interface{}) {
   mode.pool.Give(value.([]byte))
}

type BufferBytePool struct {
	pool types.BufferPoolMode
	clean   [][]byte
}

func NewBufferBytePool() types.BufferPoolMode {
	return &BufferBytePool {
		pool: BufferGetPool(BufferByteCtx{}),
	}
}

func (pool *BufferBytePool) Take(i interface{}) interface{} {
	 buf := pool.pool.Take(i.(int)).([]byte)
	 pool.clean = append(pool.clean, buf)
	 return buf
}

func (pool *BufferBytePool) Give(interface{}) {
	for _, value := range pool.clean {
		pool.pool.Give(value)
	}
}