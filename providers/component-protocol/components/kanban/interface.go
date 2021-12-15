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
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// IKanban is user-level interface for kanban.
type IKanban interface {
	cptype.IComponent
	IKanbanStdOps
}

// IKanbanStdOps defines all standard operations of kanban.
type IKanbanStdOps interface {
	// kanban-level
	// RegisterBoardCreateOp will refresh all boards.
	RegisterBoardCreateOp(opData OpBoardCreate) (opFunc cptype.OperationFunc)
	// board-level
	// RegisterBoardLoadMoreOp only return specific board data, not all boards data.
	RegisterBoardLoadMoreOp(opData OpBoardLoadMore) (opFunc cptype.OperationFunc)
	// RegisterBoardUpdateOp will refresh all boards.
	RegisterBoardUpdateOp(opData OpBoardUpdate) (opFunc cptype.OperationFunc)
	// RegisterBoardDeleteOp will refresh all boards.
	RegisterBoardDeleteOp(opData OpBoardDelete) (opFunc cptype.OperationFunc)
	// card-level
	RegisterCardMoveToOp(opData OpCardMoveTo) (opFunc cptype.OperationFunc)
}
