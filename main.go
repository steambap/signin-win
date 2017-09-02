//go:generate goversioninfo -icon=renai.ico -manifest=signin-win.exe.manifest

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var apiOrigin string

var MY_FONT = Font{PointSize: 14, Family: "微软雅黑"}

var urlConfig UrlConfig

var setting *walk.IniFileSettings

func init() {
	// dev api setting
	if apiOrigin == "" {
		apiOrigin = "http://localhost:8900"
	}
	initLoc := "0"
	setting = walk.NewIniFileSettings("settings.ini")
	if err := setting.Load(); err != nil {
		fmt.Print(err)
	} else {
		savedLoc, ok := setting.Get("loc")
		if ok {
			initLoc = savedLoc
		}
	}
	urlConfig = UrlConfig{Loc: initLoc, Date: time.Now()}
	for key, value := range bucketMap {
		bucketSlice = append(bucketSlice, LocPair{key, value})
	}
	sort.Slice(bucketSlice, func(i, j int) bool {
		var key1 = bucketSlice[i].key
		var key2 = bucketSlice[j].key
		var num1, err1 = strconv.ParseInt(key1, 10, 32)
		var num2, err2 = strconv.ParseInt(key2, 10, 32)

		if err1 != nil && err2 != nil {
			return key1 < key2
		} else if err1 != nil {
			return false
		} else if err2 != nil {
			return true
		} else {
			return num1 < num2
		}
	})
}

func main() {
	var window *walk.MainWindow
	icon, err := walk.NewIconFromResourceId(6)
	if err != nil {
		log.Print(err)
		icon = nil
	}

	_, err = MainWindow{
		Title:    "心栈签到",
		AssignTo: &window,
		MinSize:  Size{Width: 540, Height: 320},
		Layout:   VBox{},
		Font:     MY_FONT,
		Icon:     icon,
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
				Layout: Grid{Columns: 3},
				Children: []Widget{
					PushButton{
						Text: "导出日志",
						OnClicked: func() {
							if cmd, err := getUrlConfigDialog(window, RequestDailyLog); err != nil {
								walk.MsgBox(window, "导出选项弹窗错误", err.Error(), walk.MsgBoxIconError)
							} else if cmd == walk.DlgCmdOK {
								logBody, err := getDailyLog(&urlConfig)
								if err != nil {
									walk.MsgBox(window, "获取远程数据错误", err.Error(), walk.MsgBoxIconError)
								} else {
									if _, err2 := runExportDailyLogDialog(window, logBody, urlConfig.Date); err2 != nil {
										walk.MsgBox(window, "导出数据窗口错误", err.Error(), walk.MsgBoxIconError)
									}
								}
							}
						},
					},
					PushButton{
						Text: "导出周数据",
						OnClicked: func() {
							if cmd, err := getUrlConfigDialog(window, RequestWeekData); err != nil {
								walk.MsgBox(window, "导出选项弹窗错误", err.Error(), walk.MsgBoxIconError)
							} else if cmd == walk.DlgCmdOK {
								logList, err := getWeekData(&urlConfig)
								if err != nil {
									walk.MsgBox(window, "获取远程数据错误", err.Error(), walk.MsgBoxIconError)
								} else {
									if _, err2 := runExportWeekDataDialog(window, logList, urlConfig.Date); err2 != nil {
										walk.MsgBox(window, "导出数据窗口错误", err.Error(), walk.MsgBoxIconError)
									}
								}
							}
						},
					},
					PushButton{
						Text: "年统计信息",
						OnClicked: func() {
							if cmd, err := getUrlConfigDialog(window, RequestYearData); err != nil {
								walk.MsgBox(window, "导出选项弹窗错误", err.Error(), walk.MsgBoxIconError)
							} else if cmd == walk.DlgCmdOK {
								yearStats, err := getYearStats(&urlConfig)
								if err != nil {
									walk.MsgBox(window, "获取远程数据错误", err.Error(), walk.MsgBoxIconError)
								} else {
									walk.MsgBox(
										window,
										"年统计信息",
										fmt.Sprintf(
											"杯数：%v\r\n人数：%v\r\n人次：%v\r\n新人：%v",
											yearStats.CupSize,
											yearStats.NumOfTime,
											yearStats.NumOfPeople,
											yearStats.NumOfNew,
										),
										walk.MsgBoxOK)
								}
							}
						},
					},
					PushButton{
						Text:       "查看/编辑全部数据",
						ColumnSpan: 3,
						Enabled:    false,
						OnClicked: func() {
							if _, err := runOverviewDialog(window); err != nil {
								walk.MsgBox(window, "查看全部数据弹窗错误", err.Error(), walk.MsgBoxIconError)
							}
						},
					},
					Composite{
						ColumnSpan: 3,
						Layout:     VBox{MarginsZero: true},
						Children: []Widget{
							GroupBox{
								Title:  "手机签到：",
								Layout: HBox{},
								Children: []Widget{
									Label{
										Text: apiOrigin,
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run()

	if err != nil {
		log.Fatalf("Fail to Create Window:\n %v", err)
	}

	err = setting.Put("loc", urlConfig.Loc)
	if err == nil {
		setting.Save()
	} else {
		fmt.Print(err)
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

	return logBody, err
}

func getWeekData(urlConfig *UrlConfig) ([]Body, error) {
	resp, err := http.Get(apiOrigin + urlConfig.ToWeekDataUrl())
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
	var logList = make([]Body, 0, 7)
	err = json.Unmarshal(resBody, &logList)

	return logList, err
}

func getYearStats(urlConfig *UrlConfig) (*YearStats, error) {
	resp, err := http.Get(apiOrigin + urlConfig.ToYearStatsUrl())
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
	var yearStats = &YearStats{}
	err = json.Unmarshal(resBody, yearStats)
	return yearStats, err
}
