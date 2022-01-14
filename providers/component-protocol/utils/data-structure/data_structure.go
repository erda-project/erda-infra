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

package structure

const (
	String      Type = "string"
	Number      Type = "number"
	Capacity    Type = "capacity"
	TrafficRate Type = "trafficRate"
	Storage     Type = "storage"
	Timestamp   Type = "timestamp"
	Rate        Type = "rate"
	Date        Type = "date"
	Percent     Type = "percent"
	Throughput  Type = "throughput"
)

const (
	K           Precision = "K"
	M           Precision = "M"
	B           Precision = "B"
	KB          Precision = "KB"
	MB          Precision = "MB"
	GB          Precision = "GB"
	TB          Precision = "TB"
	PB          Precision = "PB"
	EB          Precision = "EB"
	ZB          Precision = "ZB"
	YB          Precision = "YB"
	BSlashS     Precision = "B/s"
	KBSlashS    Precision = "KB/s"
	MBSlashS    Precision = "MB/s"
	GBSlashS    Precision = "GB/s"
	TBSlashS    Precision = "TB/s"
	PBSlashS    Precision = "PB/s"
	EBSlashS    Precision = "EB/s"
	ZBSlashS    Precision = "ZB/s"
	YBSlashS    Precision = "YB/s"
	Nanosecond  Precision = "ns"
	Microsecond Precision = "μs"
	Millisecond Precision = "ms"
	Second      Precision = "s"
	ReqSlashS   Precision = "req/s"
)

type (
	Type          string
	Precision     string
	DataStructure struct {
		Type      Type      `json:"type"`
		Precision Precision `json:"precision"`
		Enable    bool      `json:"enable"`
	}
)
