package main

import (
	"github.com/lxn/walk"
	"time"
	"strconv"
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
