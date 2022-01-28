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

package parallel

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	type args struct {
		ctx  context.Context
		fn   func(ctx context.Context) (interface{}, error)
		opts []RunOption
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr interface{}
	}{
		{
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) (interface{}, error) {
					return 1, nil
				},
			},
			want:    1,
			wantErr: nil,
		},
		{
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) (interface{}, error) {
					return nil, context.Canceled
				},
			},
			want:    nil,
			wantErr: context.Canceled,
		},
		{
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) (interface{}, error) {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				},
				opts: []RunOption{WithTimeout(1 * time.Second)},
			},
			want:    nil,
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			future := Go(tt.args.ctx, tt.args.fn, tt.args.opts...)

			if data, err := future.Get(); !reflect.DeepEqual(data, tt.want) || !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Go() And Get() = (%v, %v), want (%v, %v)", data, err, tt.want, tt.wantErr)
			}
		})
	}
}
