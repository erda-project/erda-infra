// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package pbutil

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetValue(v interface{}) (*structpb.Value, error) {
	vv, err := structpb.NewValue(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get protobuf value from interface, err: %v", err)
	}
	return vv, nil
}

func MustGetValue(v interface{}) *structpb.Value {
	vv, err := GetValue(v)
	if err != nil {
		logrus.Warn(err)
	}
	return vv
}
