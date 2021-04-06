// Copyright 2021 Terminus
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
	"fmt"
)

type testWriter struct {
	capacity int
}

func (w *testWriter) Write(data interface{}) error {
	fmt.Printf("Write data ok")
	return nil
}

func (w *testWriter) WriteN(data ...interface{}) (int, error) {
	if w.capacity <= 0 {
		err := fmt.Errorf("buffer max capacity")
		return 0, err
	}
	w.capacity -= len(data)
	fmt.Printf("WriteN size: %d, data: %v\n", len(data), data)
	return len(data), nil
}

func (w *testWriter) Close() error {
	fmt.Println("Close")
	return nil
}

func ExampleBuffer() {
	buf := NewBuffer(&testWriter{2}, 4)
	n, err := buf.WriteN(1, 2, 3, 4, 5, 6, 7, 8, 9)
	fmt.Printf("writes: %d, buffered: %v, err: %s\n", n, buf.buf, err)

	err = buf.Close()
	fmt.Println(err)

	// Output:
	// WriteN size: 4, data: [1 2 3 4]
	// writes: 8, buffered: [5 6 7 8], err: buffer max capacity
	// Close
	// buffer max capacity
}

func ExampleBuffer_write() {
	buf := NewBuffer(&testWriter{10}, 3)
	n, err := buf.WriteN(1, 2, 3, 4, 5)
	fmt.Printf("writes: %d, buffered: %v, err: %v\n", n, buf.buf, err)

	n, err = buf.WriteN(6)
	fmt.Printf("writes: %d, buffered: %v, err: %v\n", n, buf.buf, err)

	n, err = buf.WriteN(7, 8, 9)
	fmt.Printf("writes: %d, buffered: %v, err: %v\n", n, buf.buf, err)

	n, err = buf.WriteN(10, 11)
	fmt.Printf("writes: %d, buffered: %v, err: %v\n", n, buf.buf, err)

	n, err = buf.WriteN(10, 11)
	fmt.Printf("writes: %d, buffered: %v, err: %v\n", n, buf.buf, err)

	// Output:
	// WriteN size: 3, data: [1 2 3]
	// writes: 5, buffered: [4 5], err: <nil>
	// WriteN size: 3, data: [4 5 6]
	// writes: 1, buffered: [], err: <nil>
	// writes: 3, buffered: [7 8 9], err: <nil>
	// WriteN size: 3, data: [7 8 9]
	// writes: 2, buffered: [10 11], err: <nil>
	// WriteN size: 3, data: [10 11 10]
	// writes: 2, buffered: [11], err: <nil>
}

func ExampleBuffer_for() {
	buf := NewBuffer(&testWriter{100}, 3)
	data := make([]interface{}, 10, 10)
	for i := range data {
		data[i] = i
	}
	for i := 0; i < 10; i++ {
		n, err := buf.WriteN(data[0 : 1+i%10]...)
		if err != nil {
			fmt.Println(n, buf.buf, err)
			break
		}
	}
	err := buf.Close()
	fmt.Println(buf.buf, err)

	// Output:
	// WriteN size: 3, data: [0 0 1]
	// WriteN size: 3, data: [0 1 2]
	// WriteN size: 3, data: [0 1 2]
	// WriteN size: 3, data: [3 0 1]
	// WriteN size: 3, data: [2 3 4]
	// WriteN size: 3, data: [0 1 2]
	// WriteN size: 3, data: [3 4 5]
	// WriteN size: 3, data: [0 1 2]
	// WriteN size: 3, data: [3 4 5]
	// WriteN size: 3, data: [6 0 1]
	// WriteN size: 3, data: [2 3 4]
	// WriteN size: 3, data: [5 6 7]
	// WriteN size: 3, data: [0 1 2]
	// WriteN size: 3, data: [3 4 5]
	// WriteN size: 3, data: [6 7 8]
	// WriteN size: 3, data: [0 1 2]
	// WriteN size: 3, data: [3 4 5]
	// WriteN size: 3, data: [6 7 8]
	// WriteN size: 1, data: [9]
	// Close
	// [] <nil>
}
