package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func getUrlConfigDialog(parent walk.Form, reqType RequestType) (int, error) {
	var dialog *walk.Dialog
	var dateEdit *walk.DateEdit
	listHandle := &LocListAdapter{model: &LocListModel{items: bucketSlice}}
	var label *walk.Label

	return Dialog{
		AssignTo: &dialog,
		Title:    "获取数据",
		MinSize:  Size{Width: 540, Height: 360},
		Layout:   VBox{},
		Font:     MY_FONT,
		Children: []Widget{
			Label{Text: "第一步：选择一个心栈"},
			ComboBox{
				AssignTo:     &listHandle.view,
				Model:        listHandle.model,
				CurrentIndex: locIndexOf(bucketSlice, urlConfig.Loc),
				OnCurrentIndexChanged: func() {
					loc := listHandle.model.items[listHandle.view.CurrentIndex()]
					urlConfig.Loc = loc.key
					label.SetText(urlConfig.Explain(reqType))
				},
			},
			Label{Text: "第二步：选择一个日期"},
			DateEdit{
				Date:     urlConfig.Date,
				AssignTo: &dateEdit,
				OnDateChanged: func() {
					urlConfig.Date = dateEdit.Date()
					label.SetText(urlConfig.Explain(reqType))
				},
			},
			Label{
				Text: urlConfig.Explain(reqType),
				AssignTo: &label,
			},
			PushButton{
				Text: "确认",
				OnClicked: func() {
					dialog.Accept()
				},
			},
		},
	}.Run(parent)
}
