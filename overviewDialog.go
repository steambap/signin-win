package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runOverviewDialog(parent walk.Form) (int, error) {
	var dialog *walk.Dialog
	var treeView *walk.TreeView


	return Dialog{
		AssignTo:&dialog,
		Title:"数据总览",
		MinSize:Size{Width:640, Height:480},
		Font:MY_FONT,
		Layout:VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSplitter{
						Children: []Widget{
							TreeView{
								AssignTo: &treeView,
							},
							Label{Text: "fixme"},
						},
					},
				},
			},
		},
	}.Run(parent)
}