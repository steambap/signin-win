package main

import (
	"github.com/lxn/walk"
	"time"
	"strconv"
	"fmt"
	"strings"
)

type BaseListAdapter struct {
	view  *walk.ListBox
	model *ListAdapterModel
}

// interface ListModel
type ListAdapterModel struct {
	walk.ListModelBase
	items []string
}

func (model *ListAdapterModel) ItemCount() int {
	return len(model.items)
}

func (model *ListAdapterModel) Value(index int) interface{} {
	return model.items[index]
}

type LocListAdapter struct {
	view  *walk.ComboBox
	model *LocListModel
}

type LocListModel struct {
	walk.ListModelBase
	items []LocPair
}

func (model *LocListModel) ItemCount() int {
	return len(model.items)
}

func (model *LocListModel) Value(index int) interface{} {
	return model.items[index].value
}

type UrlConfig struct {
	Loc  string
	Date time.Time
}

func (urlConfig *UrlConfig) ToDailyUrl() string {
	return "/log?date=" + urlConfig.Date.Format("2006-01-02") + "&loc=" + urlConfig.Loc
}

type Body struct {
	Names   []string `json:"names" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	Comment string   `json:"comment"`
	CupSize int      `json:"cup_size"`
}

func (log *Body) getCupSizeText() string {
	return "杯数" + strconv.FormatInt(int64(log.CupSize), 10)
}

func (log *Body) getCountText() string {
	return "人数" + strconv.FormatInt(int64(len(log.Names)), 10)
}

func (log *Body) getExportLineArr() []string {
	var cupSizeText = log.getCupSizeText()
	var countText = log.getCountText()

	return append([]string{cupSizeText, countText}, log.Names...)
}

func (log *Body) getPreview(t time.Time) string {
	tagToName := log.remixTagTable()
	ret := fmt.Sprintf("奉粥日期：%v年%v月%v日（%v）\r\n", t.Year(), int(t.Month()), t.Day(), t.Weekday().String())
	ret += fmt.Sprintf("日负责人：%v\r\n", tagToName["负责人"])
	ret += fmt.Sprintf("签到：%v\r\n", tagToName["签到"])
	ret += fmt.Sprintf("熬粥：%v\r\n", tagToName["熬粥"])
	ret += fmt.Sprintf("前行：%v\r\n", tagToName["前行"])
	ret += fmt.Sprintf("杯数：%v\r\n", log.CupSize)
	ret += fmt.Sprintf("人数：%v\r\n", len(log.Names))
	ret += fmt.Sprintf("新人：%v\r\n", tagToName["新人"])
	ret += fmt.Sprintf("摄影：%v\r\n", tagToName["摄影"])
	ret += fmt.Sprintf("日志：无\r\n")
	ret += fmt.Sprintf("文宣：\r\n")
	ret += fmt.Sprintf("结行：%v\r\n", tagToName["结行"])
	ret += fmt.Sprintf("后勤：%v\r\n", tagToName["后勤"])
	ret += fmt.Sprintf("环保：%v\r\n", tagToName["环保"])
	ret += fmt.Sprintf("奉粥：%v\r\n", tagToName["奉粥"])

	return ret
}

func (log *Body) remixTagTable() map[string]string {
	tagToName := make(map[string]string, 10)
	for index, tags := range log.Tags {
		name := log.Names[index]
		// 所有没有参加环保的都是奉粥人员
		ok := strings.Contains(tags, "环保")
		if !ok {
			names, exist := tagToName["奉粥"]
			if exist {
				tagToName["奉粥"] = names + "、" + name
			} else {
				tagToName["奉粥"] = name
			}
		}
		tagList := strings.Split(tags, "|")
		for _, tag := range tagList {
			names, exist := tagToName[tag]
			if exist {
				tagToName[tag] = names  + "、" + name
			} else {
				tagToName[tag] = name
			}
		}
	}

	return tagToName
}
