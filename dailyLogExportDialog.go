package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/xuri/excelize"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

func runExportDailyLogDialog(parent walk.Form, logBody *Body, t time.Time) (int, error) {
	var dialog *walk.Dialog
	listHandle := &BaseListAdapter{model: &ListAdapterModel{items: logBody.Names}}

	return Dialog{
		AssignTo: &dialog,
		Title:    "导出数据预览",
		Layout:   VBox{},
		MinSize:  Size{Width: 480, Height: 640},
		Font:     MY_FONT,
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					{
						Title:  "志愿者列表",
						Layout: Grid{Columns: 2},
						Children: []Widget{
							Label{Text: logBody.getCupSizeText()},
							Label{Text: logBody.getCountText()},
							ListBox{
								AssignTo:   &listHandle.view,
								Model:      listHandle.model,
								ColumnSpan: 2,
							},
							PushButton{
								Text: "导出txt格式",
								OnClicked: func() {
									fDialog := walk.FileDialog{Filter: ".txt"}
									ok, err := fDialog.ShowSave(dialog)
									if err != nil {
										walk.MsgBox(dialog, "导出txt弹窗错误", err.Error(), walk.MsgBoxIconError)
									} else if ok {
										saveDailyLogText(fDialog.FilePath, logBody)
									}
								},
							},
							PushButton{
								Text: "导出Excel格式",
								OnClicked: func() {
									fDialog := walk.FileDialog{Filter: ".xlsx"}
									ok, err := fDialog.ShowSave(dialog)
									if err != nil {
										walk.MsgBox(dialog, "导出Excel弹窗错误", err.Error(), walk.MsgBoxIconError)
									} else if ok {
										saveDailyLogExcel(fDialog.FilePath, logBody)
									}
								},
							},
						},
					},
					{
						Title:  "日志预览",
						Layout: VBox{},
						Children: []Widget{
							TextEdit{
								Text:     logBody.getPreview(t),
								ReadOnly: true,
								VScroll:  true,
							},
							PushButton{
								Text: "复制内容",
								OnClicked: func() {
									clipBoard := walk.Clipboard()
									clipBoard.Clear()
									clipBoard.SetText(logBody.getPreview(t))
								},
							},
						},
					},
				},
			},
		},
	}.Run(parent)
}

func saveDailyLogText(fsPath string, log *Body) error {
	if filepath.Ext(fsPath) != ".txt" {
		fsPath = fsPath + ".txt"
	}
	return ioutil.WriteFile(fsPath, []byte(strings.Join(log.getExportLineArr(), "\r\n")), 0644)
}

func saveDailyLogExcel(fsPath string, log *Body) error {
	if filepath.Ext(fsPath) != ".xlsx" {
		fsPath = fsPath + ".xlsx"
	}

	xlsx := excelize.NewFile()
	lines := log.getExportLineArr()
	for i, line := range lines {
		xlsx.SetCellStr("Sheet1", fmt.Sprintf("A%v", i+1), line)
	}

	return xlsx.SaveAs(fsPath)
}
