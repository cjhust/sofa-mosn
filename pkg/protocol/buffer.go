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

package protocol

import (
	"github.com/alipay/sofa-mosn/pkg/types"
	"github.com/alipay/sofa-mosn/pkg/network/buffer"
	"context"
)

var defaultDataSize = 1 << 10
var defaultHeaderSize = 1 << 5


type ProtocolBufferCtx struct{}

func (ctx ProtocolBufferCtx) Name() string {
	return "ProtocolBuffers"
}

func (ctx ProtocolBufferCtx) Init() types.BufferPoolMode {
	return nil
}

func (ctx ProtocolBufferCtx) New(interface{}) interface{} {
	pbuf := new(ProtocolBuffers)

	pbuf.bytePool = buffer.NewBufferBytePool()

	pbuf.reqData = buffer.NewIoBuffer(0)
	pbuf.reqHeader = buffer.NewIoBuffer(0)
	pbuf.rspData = buffer.NewIoBuffer(0)
	pbuf.rspHeader = buffer.NewIoBuffer(0)

	pbuf.reqHeaders = make(map[string]string)
	pbuf.reqTrailers = make(map[string]string)
	pbuf.rspHeaders = make(map[string]string)
	pbuf.rspTrailers = make(map[string]string)
	return pbuf
}

func (ctx ProtocolBufferCtx) Reset(i interface{}) {
	pbuf, _ := i.(*ProtocolBuffers)
	pbuf.bytePool.Give(nil)
	pbuf.reqData.Free()
	pbuf.rspData.Free()
	pbuf.reqHeader.Free()
	pbuf.rspHeader.Free()

	for k := range pbuf.reqHeaders {
		delete(pbuf.reqHeaders, k)
	}
	for k := range pbuf.reqTrailers {
		delete(pbuf.reqTrailers, k)
	}
	for k := range pbuf.rspHeaders {
		delete(pbuf.rspHeaders, k)
	}
	for k := range pbuf.rspTrailers {
		delete(pbuf.rspTrailers, k)
	}
}

type ProtocolBuffers struct {
	reqData     types.IoBuffer
	reqHeader   types.IoBuffer
	reqHeaders  map[string]string
	reqTrailers map[string]string

	rspData     types.IoBuffer
	rspHeader   types.IoBuffer
	rspHeaders  map[string]string
	rspTrailers map[string]string

	bytePool    types.BufferPoolMode

}

func (pbuf *ProtocolBuffers) GetByte(size int) []byte {
    return pbuf.bytePool.Take(size).([]byte)
}

func (pbuf *ProtocolBuffers) GetReqData(size int) types.IoBuffer {
	if size == 0 {
		size = defaultDataSize
	}
	pbuf.reqData.Alloc(size)
	return pbuf.reqData
}

func (pbuf *ProtocolBuffers) GetReqHeader(size int) types.IoBuffer {
	if size == 0 {
		size = defaultHeaderSize
	}
	pbuf.reqHeader.Alloc(size)
	return pbuf.reqHeader
}

func (pbuf *ProtocolBuffers) GetReqHeaders() map[string]string {
	return pbuf.reqHeaders
}

func (pbuf *ProtocolBuffers) GetReqTailers() map[string]string {
	return pbuf.reqTrailers
}

func (pbuf *ProtocolBuffers) GetRspData(size int) types.IoBuffer {
	if size == 0 {
		size = defaultDataSize
	}
	pbuf.rspData.Alloc(size)
	return pbuf.rspData
}

func (pbuf *ProtocolBuffers) GetRspHeader(size int) types.IoBuffer {
	if size == 0 {
		size = defaultHeaderSize
	}
	pbuf.rspHeader.Alloc(size)
	return pbuf.rspHeader
}

func (pbuf *ProtocolBuffers) GetRspHeaders() map[string]string {
	return pbuf.rspHeaders
}

func (pbuf *ProtocolBuffers) GetRspTailers() map[string]string {
	return pbuf.rspTrailers
}

func ProtocolBuffersByContent(context context.Context) *ProtocolBuffers {
	ctx := buffer.BufferPoolContext(context)
	if ctx == nil {
		return nil
	}
	return ctx.Find(ProtocolBufferCtx{}, nil).(*ProtocolBuffers)
}

