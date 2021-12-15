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
)

type provider struct {
	impl.DefaultKanban
}

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
		},
	}
)

// RegisterInitializeOp .
func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		data := kanban.Data{
			Boards:     []kanban.Board{boardUrgent, boardNormal},
			Operations: nil,
		}
		boardUrgent.Cards = append([]kanban.Card{}, card1)
		boardUrgent.Total = 1
		p.StdDataPtr = &data
	}
}

// RegisterRenderingOp .
func (p *provider) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return p.RegisterInitializeOp()
}

// RegisterCardMoveToOp .
func (p *provider) RegisterCardMoveToOp(opData kanban.OpCardMoveTo) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		fmt.Println("hello world, i am drag:", opData)
		switch v := (cptype.ExtraMap)(*p.StdStatePtr).String("DropTarget"); v {
		case boardLabelUrgent:
			boardUrgent.Cards = append(boardUrgent.Cards, card1)
			boardNormal.Cards = nil
			p.StdDataPtr.Boards = []kanban.Board{boardUrgent, boardNormal}
		case boardLabelNormal:
			boardUrgent.Cards = nil
			boardNormal.Cards = append(boardNormal.Cards, card1)
			p.StdDataPtr.Boards = []kanban.Board{boardUrgent, boardNormal}
		default:
			panic(fmt.Errorf("invalid drop target: %s", v))
		}
	}
}

// RegisterBoardLoadMoreOp .
func (p *provider) RegisterBoardLoadMoreOp(opData kanban.OpBoardLoadMore) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		fmt.Println("hello change page no op:", opData)
	}
}

// RegisterBoardCreateOp .
func (p *provider) RegisterBoardCreateOp(opData kanban.OpBoardCreate) (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) {
		fmt.Println("hello create board op:", opData)
		p.StdDataPtr.Boards = append(p.StdDataPtr.Boards, kanban.Board{ID: opData.ClientData.Title, Title: opData.ClientData.Title})
	}
}

// Initialize .
func (p *provider) Initialize(sdk *cptype.SDK) { return }

// Visible .
func (p *provider) Visible(sdk *cptype.SDK) bool { return true }

// RegisterBoardUpdateOp .
func (p *provider) RegisterBoardUpdateOp(opData kanban.OpBoardUpdate) (opFunc cptype.OperationFunc) {
	return nil
}

// RegisterBoardDeleteOp .
func (p *provider) RegisterBoardDeleteOp(opData kanban.OpBoardDelete) (opFunc cptype.OperationFunc) {
	return nil
}
