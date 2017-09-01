package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runOverviewDialog(parent walk.Form) (int, error) {
	var dialog *walk.Dialog
	var treeView *walk.TreeView
	treeModel := newLocTreeModel()

	return Dialog{
		AssignTo: &dialog,
		Title:    "数据总览",
		MinSize:  Size{Width: 800, Height: 640},
		Font:     MY_FONT,
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSplitter{
						Children: []Widget{
							ListBox{},
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
