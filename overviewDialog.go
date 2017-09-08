package main

import (
	"encoding/json"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"net/http"
)

func runOverviewDialog(parent walk.Form) (int, error) {
	var dialog *walk.Dialog
	listHandle := &LocListBoxAdapter{model: &LocListModel{items: bucketSlice}}
	var treeView *walk.TreeView
	treeModel := newLocTreeModel()

	return Dialog{
		AssignTo: &dialog,
		Title:    "数据总览",
		MinSize:  Size{Width: 1000, Height: 600},
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
								OnCurrentIndexChanged: func() {
									locPair := listHandle.model.items[listHandle.view.CurrentIndex()]
									keys, err := scanBucket(locPair.key)
									if err != nil {
										walk.MsgBox(dialog, "获取远程数据失败", err.Error(), walk.MsgBoxOK)
									} else {
										log.Print(keys)
									}
								},
							},
							TreeView{
								AssignTo: &treeView,
								Model:    treeModel,
							},
							Label{Text: "fixme"},
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
