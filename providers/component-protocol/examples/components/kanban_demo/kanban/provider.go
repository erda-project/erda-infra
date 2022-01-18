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

package kanban

import (
	"fmt"

	"github.com/erda-project/erda-infra/providers/component-protocol/components/kanban"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/kanban/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
	"github.com/erda-project/erda-infra/providers/i18n"
)

type component struct {
	impl.DefaultKanban

	Tran i18n.Translator `translator:"notbelong"`

	// custom type
	StatePtr    *CustomState
	InParamsPtr *CustomInParams
}

const (
	nonStdOp1Key = "NonStdOp1"
)

const (
	boardLabelUrgent = "URGENT"
	boardLabelNormal = "NORMAL"
)

var (
	boardUrgent = kanban.Board{
		ID:       boardLabelUrgent,
		Title:    "紧急",
		Cards:    nil,
		PageNo:   1,
		PageSize: 20,
		Total:    0,
	}
	boardNormal = kanban.Board{
		ID:       boardLabelNormal,
		Title:    "普通",
		Cards:    nil,
		PageNo:   1,
		PageSize: 20,
		Total:    0,
	}
	card1 = kanban.Card{
		ID:    "1",
		Title: "task1",
		Extra: cptype.Extra{
			Extra: map[string]interface{}{
				"userID": "2",
				"status": "Done",
			},
		},
		Operations: map[cptype.OperationKey]cptype.Operation{
			kanban.OpCardMoveTo{}.OpKey(): cputil.NewOpBuilder().
				WithAsync(true).
				WithServerDataPtr(&kanban.OpCardMoveToServerData{AllowedTargetBoardIDs: []string{boardLabelUrgent, boardLabelNormal}}).
				Build(),
			// Test Non-Std op
			nonStdOp1Key: cputil.NewOpBuilder().Build(),
		},
	}
)

func (c *component) registerNonStdOp1() cptype.OperationFunc {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		fmt.Println("This is NonStdOp1")
		return nil
	}
}

// RegisterCompNonStdOps .
func (c *component) RegisterCompNonStdOps() (opFuncs map[cptype.OperationKey]cptype.OperationFunc) {
	return map[cptype.OperationKey]cptype.OperationFunc{
		nonStdOp1Key: c.registerNonStdOp1(),
	}
}

// RegisterInitializeOp .
func (c *component) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		data := kanban.Data{
			Boards:     []kanban.Board{boardUrgent, boardNormal},
			Operations: nil,
		}
		// multi instance demo
		switch sdk.Comp.Name {
		case "instance1":
			fmt.Println("this is instance1")
			c.StatePtr.Name = "instance1"
		case "instance2":
			fmt.Println("this is instance2")
			c.StatePtr.Name = "instance2"
		default:
			fmt.Println("this is a kanban instance")
			c.StatePtr.Name = "Bob"
		}
		boardUrgent.Cards = append([]kanban.Card{}, card1)
		boardUrgent.Total = 1
		c.StdDataPtr = &data
		// custom inParams
		c.InParamsPtr.ProjectID = 20
		fmt.Println("hello", c.Tran.Text(i18n.LanguageCodes{{Code: "zh"}}, "hello"))
		return nil
	}
}

// RegisterRenderingOp .
func (c *component) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return c.RegisterInitializeOp()
}

// RegisterCardMoveToOp .
func (c *component) RegisterCardMoveToOp(opData kanban.OpCardMoveTo) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		fmt.Println("hello world, i am drag:", opData)
		switch v := (cptype.ExtraMap)(*c.StdStatePtr).String("DropTarget"); v {
		case boardLabelUrgent:
			boardUrgent.Cards = append(boardUrgent.Cards, card1)
			boardNormal.Cards = nil
			c.StdDataPtr.Boards = []kanban.Board{boardUrgent, boardNormal}
		case boardLabelNormal:
			boardUrgent.Cards = nil
			boardNormal.Cards = append(boardNormal.Cards, card1)
			c.StdDataPtr.Boards = []kanban.Board{boardUrgent, boardNormal}
		default:
			panic(fmt.Errorf("invalid drop target: %s", v))
		}
		return nil
	}
}

// RegisterBoardLoadMoreOp .
func (c *component) RegisterBoardLoadMoreOp(opData kanban.OpBoardLoadMore) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		fmt.Println("hello change page no op:", opData)
		return nil
	}
}

// RegisterBoardCreateOp .
func (c *component) RegisterBoardCreateOp(opData kanban.OpBoardCreate) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		fmt.Println("hello create board op:", opData)
		c.StdDataPtr.Boards = append(c.StdDataPtr.Boards, kanban.Board{ID: opData.ClientData.Title, Title: opData.ClientData.Title})
		return nil
	}
}

// Initialize .
func (c *component) Initialize(sdk *cptype.SDK) { return }

// Visible .
func (c *component) Visible(sdk *cptype.SDK) bool { return true }

// RegisterBoardUpdateOp .
func (c *component) RegisterBoardUpdateOp(opData kanban.OpBoardUpdate) (opFunc cptype.OperationFunc) {
	return nil
}

// RegisterBoardDeleteOp .
func (c *component) RegisterBoardDeleteOp(opData kanban.OpBoardDelete) (opFunc cptype.OperationFunc) {
	return nil
}
