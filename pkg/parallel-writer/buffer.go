// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writer

import (
	"github.com/recallsong/go-utils/errorx"
)

// Buffer .
type Buffer struct {
	w       Writer
	buf     []interface{}
	maxSize int
}

// NewBuffer .
func NewBuffer(w Writer, max int) *Buffer {
	return &Buffer{
		w:       w,
		buf:     make([]interface{}, 0, max),
		maxSize: max,
	}
}

// Write .
func (b *Buffer) Write(data interface{}) error {
	if len(b.buf)+1 > b.maxSize {
		err := b.Flush()
		if err != nil {
			return err
		}
	}
	b.buf = append(b.buf, data)
	return nil
}

// WriteN returns the number of buffers written to the data.
// if a Flush error occurs, the error will be returned
func (b *Buffer) WriteN(data ...interface{}) (int, error) {
	alen := len(b.buf)
	blen := len(data)
	if alen+blen < b.maxSize {
		b.buf = append(b.buf, data...)
		return blen, nil
	}
	writes := 0
	if alen >= b.maxSize {
		// never reached
		err := b.Flush()
		if err != nil {
			return 0, nil
		}
	} else if alen > 0 {
		writes = b.maxSize - alen
		b.buf = append(b.buf, data[0:writes]...)
		err := b.Flush()
		if err != nil {
			return writes, err
		}
		data = data[writes:]
		blen -= writes
	}
	for blen > b.maxSize {
		b.buf = append(b.buf, data[0:b.maxSize]...)
		writes += b.maxSize
		err := b.Flush()
		if err != nil {
			return writes, err
		}
		data = data[b.maxSize:]
		blen -= b.maxSize
	}
	if blen > 0 {
		b.buf = append(b.buf, data...)
		writes += blen
	}
	return writes, nil
}

// Flush .
func (b *Buffer) Flush() error {
	l := len(b.buf)
	if l > 0 {
		n, err := b.w.WriteN(b.buf...)
		b.buf = b.buf[0 : l-n]
		if err != nil {
			return err
		}
	}
	return nil
}

// Size .
func (b *Buffer) Size() int {
	return len(b.buf)
}

// Data .
func (b *Buffer) Data() []interface{} {
	return b.buf
}

// Close .
func (b *Buffer) Close() error {
	return errorx.NewMultiError(b.Flush(), b.w.Close()).MaybeUnwrap()
}
