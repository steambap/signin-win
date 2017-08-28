package main

import (
	"encoding/json"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const apiOrigin = "http://localhost:8900"

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
								OnClicked: func() {
									if cmd, err := runDailyLogDialog(window, &urlConfig); err != nil {
										walk.MsgBox(window, "导出选项弹窗错误", err.Error(), walk.MsgBoxIconError)
									} else if cmd == walk.DlgCmdOK {
										logBody, err := getDailyLog(&urlConfig)
										if err != nil {
											walk.MsgBox(window, "获取远程数据错误", err.Error(), walk.MsgBoxIconError)
										} else {
											if _, err2 := runExportDailyLogDialog(window, logBody); err2 != nil {
												walk.MsgBox(window, "导出数据窗口错误", err.Error(), walk.MsgBoxIconError)
											}
										}
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

func getDailyLog(urlConfig *UrlConfig) (*Body, error) {
	resp, err := http.Get(apiOrigin + urlConfig.ToDailyUrl())
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
	var logBody = &Body{
		Names:   make([]string, 0),
		Tags:    make([]string, 0),
		Comment: "",
		CupSize: -1,
	}
	err = json.Unmarshal(resBody, logBody)
	if err != nil {
		return nil, err
	}

	return logBody, nil
}
