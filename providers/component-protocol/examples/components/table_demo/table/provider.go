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

package table

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/components/commodel"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/table"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/table/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

type provider struct {
	impl.DefaultTable
}

const (
	columnMergedTitle    table.ColumnKey = "mergedTitle"
	columnIcon           table.ColumnKey = "icon"
	columnTitle          table.ColumnKey = "title"
	columnLabels         table.ColumnKey = "labels"
	columnPriority       table.ColumnKey = "priority"
	columnState          table.ColumnKey = "state"
	columnAssignee       table.ColumnKey = "assignee"
	columnFinishedAt     table.ColumnKey = "finishedAt"
	columnMoreOperations table.ColumnKey = "moreOperations"
	columnProgress       table.ColumnKey = "progress"
)

func (p *provider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		p.StdDataPtr = &table.Data{
			Table: table.Table{
				Columns: table.ColumnsInfo{
					Merges: map[table.ColumnKey]table.MergedColumn{
						columnMergedTitle: {Orders: []table.ColumnKey{columnIcon, columnTitle, columnLabels}},
					},
					Orders: []table.ColumnKey{columnTitle, columnPriority, columnState, columnAssignee, columnFinishedAt, columnProgress, columnMoreOperations},
					ColumnsMap: map[table.ColumnKey]table.Column{
						columnTitle:      {Title: "标题"},
						columnPriority:   {Title: "优先级"},
						columnState:      {Title: "状态"},
						columnAssignee:   {Title: "处理人"},
						columnFinishedAt: {Title: "截止日期"},
						columnProgress:   {Title: "进度条"},
					},
				},
				Rows: []table.Row{
					{
						ID:         "issue-id-1",
						Selectable: true,
						Selected:   false,
						CellsMap: map[table.ColumnKey]table.Cell{
							columnIcon:  table.NewIconCell(*commodel.NewTypedIcon("ISSUE_ICON.issue.TASK")).Build(),
							columnTitle: table.NewTextCell("【服务监控】增加链路查询页面").Build(),
							columnLabels: table.NewLabelsCell(
								commodel.Labels{
									Labels: []commodel.Label{
										{
											ID:    "label-id-1",
											Title: "area/监控",
											Color: commodel.ColorRed,
										},
										{
											ID:    "label-id-2",
											Title: "team/前端",
											Color: commodel.ColorPurple,
										},
									},
								},
							).Build(),
							columnPriority: table.NewDropDownMenuCell(commodel.DropDownMenu{
								Menus: []commodel.DropDownMenuItem{
									{
										ID:       "urgent",
										Text:     "紧急",
										Icon:     commodel.NewTypedIcon("ISSUE_ICON.priority.URGENT"),
										Selected: false,
										Disabled: false,
										Hidden:   false,
									},
									{
										ID:       "high",
										Text:     "高",
										Icon:     commodel.NewTypedIcon("ISSUE_ICON.priority.HIGH"),
										Selected: false,
										Disabled: false,
										Hidden:   false,
									},
									{
										ID:       "normal",
										Text:     "中",
										Icon:     commodel.NewTypedIcon("ISSUE_ICON.priority.NORMAL"),
										Selected: true,
										Disabled: false,
										Hidden:   false,
									},
									{
										ID:       "low",
										Text:     "低",
										Icon:     commodel.NewTypedIcon("ISSUE_ICON.priority.LOW"),
										Selected: false,
										Disabled: false,
										Hidden:   false,
									},
								},
								Operations: map[cptype.OperationKey]cptype.Operation{
									commodel.OpDropDownMenuChange{}.OpKey(): cputil.NewOpBuilder().Build(),
								},
							}).Build(),
							columnState: table.NewDropDownMenuCell(commodel.DropDownMenu{
								Menus: []commodel.DropDownMenuItem{
									{
										ID:       "state-id-for-open",
										Text:     "待处理",
										Selected: false,
										Disabled: true,
										Hidden:   true,
										Tip:      "无法转移",
									},
									{
										ID:       "state-id-for-working",
										Text:     "进行中",
										Selected: true,
										Disabled: false,
										Hidden:   false,
									},
									{
										ID:       "state-id-for-done",
										Text:     "已完成",
										Selected: false,
										Disabled: false,
										Hidden:   false,
									},
									{
										ID:       "state-id-for-abandoned",
										Text:     "已作废",
										Selected: false,
										Disabled: false,
										Hidden:   false,
									},
								},
								Operations: map[cptype.OperationKey]cptype.Operation{
									commodel.OpDropDownMenuChange{}.OpKey(): cputil.NewOpBuilder().Build(),
								},
							}).Build(),
							columnAssignee: table.NewUserSelectorCell(commodel.UserSelector{
								Scope:           "project",
								ScopeID:         "1000300",
								SelectedUserIDs: []string{"92"},
								Operations: map[cptype.OperationKey]cptype.Operation{
									commodel.OpUserSelectorChange{}.OpKey(): cputil.NewOpBuilder().Build(),
								},
							}).Build(),
							columnFinishedAt: table.NewTextCell("2021-12-29").Build(),
							columnProgress: table.NewProgressBarCell(commodel.ProgressBar{
								BarPercent: 70,
								Text:       "7/10",
								Tip:        "7 executed of total 10",
								Status:     commodel.ProcessingStatus,
							}).Build(),
						},
						Operations: map[cptype.OperationKey]cptype.Operation{
							table.OpRowSelect{}.OpKey(): cputil.NewOpBuilder().Build(),
							table.OpRowAdd{}.OpKey():    cputil.NewOpBuilder().Build(),
							table.OpRowEdit{}.OpKey():   cputil.NewOpBuilder().Build(),
							table.OpRowDelete{}.OpKey(): cputil.NewOpBuilder().Build(),
						},
					},
					{
						ID: "issue-id-2",
						CellsMap: map[table.ColumnKey]table.Cell{
							columnProgress: table.NewProgressBarCell(commodel.ProgressBar{
								BarCompletedNum: 7,
								BarTotalNum:     10,
								Status:          commodel.ProcessingStatus,
							}).Build(),
						},
					},
					{
						ID: "issue-id-3",
					},
					{
						ID: "pipeline-definition-1",
						CellsMap: map[table.ColumnKey]table.Cell{
							columnMoreOperations: table.NewMoreOperationsCell(commodel.MoreOperations{
								Ops: []commodel.MoreOpItem{
									{
										ID:   "star",
										Text: "标星",
										Operations: map[cptype.OperationKey]cptype.Operation{
											commodel.OpMoreOperationsItemClick{}.OpKey(): cputil.NewOpBuilder().
												WithServerDataPtr(&commodel.OpMoreOperationsItemClickServerData{}).
												Build(),
										},
									},
									{
										ID:   "goto-detail-page",
										Text: "查看详情",
										Operations: map[cptype.OperationKey]cptype.Operation{
											commodel.OpMoreOperationsItemClickGoto{}.OpKey(): cputil.NewOpBuilder().
												WithServerDataPtr(&commodel.OpMoreOperationsItemClickGotoServerData{
													OpClickGotoServerData: commodel.OpClickGotoServerData{
														JumpOut: false,
														Target:  "projectPipelineDetail",
														Params: map[string]interface{}{
															"pipelineDefinitionID": "1",
														},
														Query: nil,
													}}).
												Build(),
										},
									},
								},
							}).Build(),
						},
					},
				},
				PageNo:   1,
				PageSize: 10,
				Total:    1,
			},
			Operations: map[cptype.OperationKey]cptype.Operation{
				table.OpTableChangePage{}.OpKey(): cputil.NewOpBuilder().WithServerDataPtr(&table.OpTableChangePageServerData{}).Build(),
				table.OpTableChangeSort{}.OpKey(): cputil.NewOpBuilder().Build(),
				table.OpBatchRowsHandle{}.OpKey(): cputil.NewOpBuilder().WithText("批量操作").WithServerDataPtr(&table.OpBatchRowsHandleServerData{
					Options: []table.OpBatchRowsHandleOption{
						{
							ID:            "delete",
							Text:          "删除",
							AllowedRowIDs: []string{"row1", "row2"},
						},
						{
							ID:              "start",
							Text:            "启动",
							ForbiddenRowIDs: []string{"row2"},
						},
					},
				}).Build(),
			},
		}
		return nil
	}
}

func (p *provider) RegisterRenderingOp() (opFunc cptype.OperationFunc) {
	return p.RegisterInitializeOp()
}

func (p *provider) RegisterTablePagingOp(opData table.OpTableChangePage) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterTableChangePageOp(opData table.OpTableChangePage) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterTableSortOp(opData table.OpTableChangeSort) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterBatchRowsHandleOp(opData table.OpBatchRowsHandle) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowSelectOp(opData table.OpRowSelect) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowAddOp(opData table.OpRowAdd) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowEditOp(opData table.OpRowEdit) (opFunc cptype.OperationFunc) {
	return nil
}

func (p *provider) RegisterRowDeleteOp(opData table.OpRowDelete) (opFunc cptype.OperationFunc) {
	return nil
}
