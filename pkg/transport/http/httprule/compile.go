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

// Reference: https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.3.0/internal/httprule/compile.go

package httprule

import (
	"github.com/erda-project/erda-infra/pkg/transport/http/utilities"
)

const (
	// SupportPackageIsVersion1 .
	SupportPackageIsVersion1 = 1
	opcodeVersion            = 1
)

// Template is a compiled representation of path templates.
type Template struct {
	// Version is the version number of the format.
	Version int
	// OpCodes is a sequence of operations.
	OpCodes []int
	// Pool is a constant pool
	Pool []string
	// Verb is a VERB part in the template.
	Verb string
	// Fields is a list of field paths bound in this template.
	Fields []string
	// Original template (example: /v1/a_bit_of_everything)
	Template string
}

// Compiler compiles utilities representation of path templates into marshallable operations.
// They can be unmarshalled by runtime.NewPattern.
type Compiler interface {
	Compile() Template
}

type op struct {
	// code is the opcode of the operation
	code utilities.OpCode

	// str is a string operand of the code.
	// num is ignored if str is not empty.
	str string

	// num is a numeric operand of the code.
	num int
}

func (w wildcard) compile() []op {
	return []op{
		{code: utilities.OpPush},
	}
}

func (w deepWildcard) compile() []op {
	return []op{
		{code: utilities.OpPushM},
	}
}

func (l literal) compile() []op {
	return []op{
		{
			code: utilities.OpLitPush,
			str:  string(l),
		},
	}
}

func (v variable) compile() []op {
	var ops []op
	for _, s := range v.segments {
		ops = append(ops, s.compile()...)
	}
	ops = append(ops, op{
		code: utilities.OpConcatN,
		num:  len(v.segments),
	}, op{
		code: utilities.OpCapture,
		str:  v.path,
	})

	return ops
}

func (t template) Compile() Template {
	var rawOps []op
	for _, s := range t.segments {
		rawOps = append(rawOps, s.compile()...)
	}

	var (
		ops    []int
		pool   []string
		fields []string
	)
	consts := make(map[string]int)
	for _, op := range rawOps {
		ops = append(ops, int(op.code))
		if op.str == "" {
			ops = append(ops, op.num)
		} else {
			if _, ok := consts[op.str]; !ok {
				consts[op.str] = len(pool)
				pool = append(pool, op.str)
			}
			ops = append(ops, consts[op.str])
		}
		if op.code == utilities.OpCapture {
			fields = append(fields, op.str)
		}
	}
	return Template{
		Version:  opcodeVersion,
		OpCodes:  ops,
		Pool:     pool,
		Verb:     t.verb,
		Fields:   fields,
		Template: t.template,
	}
}
