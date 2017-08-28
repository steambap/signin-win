package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

func runExportDailyLogDialog(parent walk.Form, logBody *Body) (int, error) {
	var dialog *walk.Dialog
	listHandle := &BaseListAdapter{model: &ListAdapterModel{items: logBody.Names}}
	log.Print(logBody)

	return Dialog{
		AssignTo: &dialog,
		Title:    "导出数据预览",
		Layout:   VBox{},
		MinSize:  Size{Width: 640, Height: 480},
		Children: []Widget{
			ListBox{
				AssignTo: &listHandle.view,
				Model:    listHandle.model,
			},
		},
	}.Run(parent)
}
