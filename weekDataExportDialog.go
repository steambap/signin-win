package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/xuri/excelize"
	"path/filepath"
	"time"
)

func runExportWeekDataDialog(parent walk.Form, logList []Body, t time.Time) (int, error) {
	var dialog *walk.Dialog
	var weekStats = calcWeekData(logList)
	var dayCountExplain string
	if len(logList) == 7 {
		dayCountExplain = "本周 7 天数据齐全"
	} else {
		dayCountExplain = fmt.Sprintf("本周数据只有 %d 天", len(logList))
	}

	return Dialog{
		AssignTo: &dialog,
		Title:    "周数据预览",
		Layout:   Grid{Columns: 2},
		Font:     MY_FONT,
		MinSize:  Size{Width: 480, Height: 320},
		Children: []Widget{
			Label{
				Text:       dayCountExplain,
				ColumnSpan: 2,
			},
			Label{
				Text: fmt.Sprintf("杯数：%v", weekStats.CupSize),
			},
			Label{
				Text: fmt.Sprintf("人次：%v", weekStats.NumOfTime),
			},
			Label{
				Text: fmt.Sprintf("人数：%v", weekStats.NumOfPeople),
			},
			Label{
				Text: fmt.Sprintf("新人：%v", weekStats.NumOfNew),
			},
			PushButton{
				Text:       "导出Excel格式",
				ColumnSpan: 2,
				OnClicked: func() {
					fDialog := walk.FileDialog{Filter: ".xlsx"}
					ok, err := fDialog.ShowSave(dialog)
					if err != nil {
						walk.MsgBox(dialog, "导出Excel弹窗错误", err.Error(), walk.MsgBoxIconError)
					} else if ok {
						saveWeekDataExcel(fDialog.FilePath, weekStats)
					}
				},
			},
		},
	}.Run(parent)
}

func saveWeekDataExcel(fsPath string, stats *WeekStats) error {
	if filepath.Ext(fsPath) != ".xlsx" {
		fsPath = fsPath + ".xlsx"
	}

	xlsx := excelize.NewFile()
	xlsx.SetCellStr("Sheet1", "A1", fmt.Sprintf("仁爱%v心栈周报", bucketMap[urlConfig.Loc]))
	xlsx.SetCellStr("Sheet1", "A2", "奉粥杯数")
	xlsx.SetCellInt("Sheet1", "A3", stats.CupSize)
	xlsx.SetCellStr("Sheet1", "B2", "志愿者人次")
	xlsx.SetCellInt("Sheet1", "B3", stats.NumOfTime)
	xlsx.SetCellStr("Sheet1", "C2", "志愿者人数")
	xlsx.SetCellInt("Sheet1", "C3", stats.NumOfPeople)
	xlsx.SetCellStr("Sheet1", "D2", "新志愿者人数")
	xlsx.SetCellInt("Sheet1", "D3", stats.NumOfNew)
	xlsx.SetCellStr("Sheet1", "E2", "接受善款总额")
	xlsx.SetCellInt("Sheet1", "E3", 0)
	xlsx.SetCellStr("Sheet1", "A4", "制作人：石威林")

	for idx, name := range stats.Names {
		xlsx.SetCellStr("Sheet1", fmt.Sprintf("C%v", idx+5), name)
	}

	return xlsx.SaveAs(fsPath)
}
