package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	//"log"
)

func runDailyLogDialog(parent walk.Form, data *UrlConfig) (int, error) {
	var dialog *walk.Dialog
	var dateEdit *walk.DateEdit
	listHandle := &LocListAdapter{model: &LocListModel{items: bucketSlice}}

	return Dialog{
		AssignTo: &dialog,
		Title:    "获取日志",
		MinSize:  Size{Width: 640, Height: 480},
		Layout:   VBox{},
		Children: []Widget{
			ComboBox{
				AssignTo: &listHandle.view,
				Model:    listHandle.model,
				CurrentIndex: locIndexOf(bucketSlice, data.Loc),
				OnCurrentIndexChanged: func() {
					loc := listHandle.model.items[listHandle.view.CurrentIndex()]
					data.Loc = loc.key
				},
			},
			DateEdit{
				Date:     data.Date,
				AssignTo: &dateEdit,
				OnDateChanged: func() {
					data.Date = dateEdit.Date()
				},
			},
			PushButton{
				Text: "OK",
				OnClicked: func() {
					dialog.Accept()
				},
			},
		},
	}.Run(parent)
}
