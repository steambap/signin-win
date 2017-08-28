package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"fmt"
	"time"
)

func init() {
	for key, value := range bucketMap {
		bucketSlice = append(bucketSlice, LocPair{key, value})
	}
}

func main() {
	var window *walk.MainWindow
	urlConfig := UrlConfig{Loc: "0", Date: time.Now()}

	_, err := MainWindow{
		Title:    "心栈签到",
		AssignTo: &window,
		MinSize:  Size{Width: 800, Height: 600},
		Layout:   VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: "文件",
				Items: []MenuItem{
					Action{
						Text: "退出",
						OnTriggered: func() {
							window.Close()
						},
					},
				},
			},
		},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{
								Text: "导出日志",
								OnClicked:func() {
									if cmd, err := runDailyLogDialog(window, &urlConfig); err != nil {
										walk.MsgBox(window, "错误", fmt.Sprintf("原因：%v", err), walk.MsgBoxIconError)
									} else if cmd == walk.DlgCmdOK {
										walk.MsgBox(window, "!", fmt.Sprint(urlConfig), walk.MsgBoxOK)
									}
								},
							},
							PushButton{
								Text:    "导出周数据",
								Enabled: false,
							},
							PushButton{
								Text:    "导出年数据",
								Enabled: false,
							},
						},
					},
					PushButton{
						Text:    "查看/编辑全部数据",
						Enabled: false,
					},
				},
			},
			Label{
				Text: "数据预览区域",
			},
		},
	}.Run()

	if err != nil {
		log.Fatalf("Fail to Create Window:\n %v", err)
	}
}
