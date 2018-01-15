package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runOverviewDialog(parent walk.Form) (int, error) {
	var dialog *walk.Dialog
	listHandle := &LocListBoxAdapter{model: &LocListModel{items: bucketSlice}}
	var treeView *walk.TreeView
	treeModel := newEmptyTreeModel()

	var numInput *walk.NumberEdit
	var txtInput *walk.TextEdit

	return Dialog{
		AssignTo: &dialog,
		Title:    "数据总览",
		MinSize:  Size{Width: 1024, Height: 768},
		Font:     MY_FONT,
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSplitter{
						Children: []Widget{
							ListBox{
								AssignTo: &listHandle.view,
								Model:    listHandle.model,
								MaxSize:  Size{Width: 200, Height: 600},
								OnCurrentIndexChanged: func() {
									locPair := listHandle.model.items[listHandle.view.CurrentIndex()]
									keys, err := scanBucket(locPair.key)
									if err != nil {
										walk.MsgBox(dialog, "获取远程数据失败", err.Error(), walk.MsgBoxOK)
									} else {
										treeView.SetModel(treeModelFromList(keys))
									}
								},
							},
							TreeView{
								AssignTo: &treeView,
								Model:    treeModel,
								MaxSize:  Size{Width: 200, Height: 600},
								OnCurrentItemChanged: func() {
									sel := treeView.CurrentItem()
									switch sel.(type) {
									case *DayItem:
										date := sel.(*DayItem).t
										loc := listHandle.model.items[listHandle.view.CurrentIndex()].key
										dailyLog, err := getDailyLog(&UrlConfig{loc, date})
										if err != nil {
											walk.MsgBox(dialog, "获取日志错误", err.Error(), walk.MsgBoxOK)
										} else {
											numInput.SetValue(float64(dailyLog.CupSize))
											namesWithTags := dailyLog.getNamesWithTags()
											txtInput.SetText(strings.Join(namesWithTags, "\r\n"))
										}
									}
								},
							},
							Composite{
								Layout: Grid{Columns: 3},
								Children: []Widget{
									Label{
										Text:       "杯数",
										ColumnSpan: 1,
										RowSpan:    1,
									},
									NumberEdit{
										AssignTo:   &numInput,
										Value:      0.0,
										ColumnSpan: 2,
										RowSpan:    1,
										ReadOnly:   true,
									},
									TextEdit{
										AssignTo:   &txtInput,
										ColumnSpan: 3,
										RowSpan:    1,
										ReadOnly:   true,
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run(parent)
}

func scanBucket(location string) ([]string, error) {
	resp, err := http.Get(apiOrigin + "/loc/" + location)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var keys = make([]string, 0)
	err = json.Unmarshal(resBody, &keys)

	return keys, err
}
