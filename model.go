package main

import (
	"fmt"
	"github.com/lxn/walk"
	"strconv"
	"strings"
	"time"
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
	return model.items[index].ToString()
}

// Location tree model implementation
// Tree Item interface
type RootLocation struct {
	LocPair
	children []*MonthItem
}

func (tree *RootLocation) Text() string {
	return "NO." + tree.key + " " + tree.value
}

func (*RootLocation) Parent() walk.TreeItem {
	return nil
}

var fakeArr = [...]string{"4", "6", "8"}

func (tree *RootLocation) ChildCount() int {
	// FIXME run HTTP request to do lazy population
	if tree.children == nil {
		children := make([]*MonthItem, 0, len(fakeArr))
		for _, num := range fakeArr {
			monthItem := &MonthItem{
				name:     "2017-0" + num,
				parent:   tree,
				children: make([]*DayItem, 0, len(fakeArr)),
			}
			for idx := range fakeArr {
				dayItem := &DayItem{
					name:   "2017-0" + num + "-0" + strconv.FormatInt(int64(idx), 10),
					parent: monthItem,
				}
				monthItem.children = append(monthItem.children, dayItem)
			}

			children = append(children, monthItem)
		}

		tree.children = children
	}

	return len(tree.children)
}

func (tree *RootLocation) ChildAt(index int) walk.TreeItem {
	return tree.children[index]
}

type MonthItem struct {
	name     string
	parent   *RootLocation
	children []*DayItem
}

func (tree *MonthItem) Text() string {
	return tree.name
}

func (tree *MonthItem) Parent() walk.TreeItem {
	return tree.parent
}

func (tree *MonthItem) ChildCount() int {
	return len(tree.children)
}

func (tree *MonthItem) ChildAt(index int) walk.TreeItem {
	return tree.children[index]
}

type DayItem struct {
	name   string
	parent *MonthItem
}

func (tree *DayItem) Text() string {
	return tree.name
}

func (tree *DayItem) Parent() walk.TreeItem {
	return tree.parent
}

func (*DayItem) ChildCount() int {
	return 0
}

func (*DayItem) ChildAt(index int) walk.TreeItem {
	return nil
}

type LocTreeModel struct {
	walk.TreeModelBase
	roots []*RootLocation
}

func newLocTreeModel() *LocTreeModel {
	roots := make([]*RootLocation, 0, len(bucketSlice))
	for _, loc := range bucketSlice {
		roots = append(roots, &RootLocation{LocPair: loc})
	}

	return &LocTreeModel{roots: roots}
}

func (*LocTreeModel) LazyPopulation() bool {
	// we do not want to scan the whole database at start
	return true
}

func (tree *LocTreeModel) RootCount() int {
	return len(tree.roots)
}

func (tree *LocTreeModel) RootAt(index int) walk.TreeItem {
	return tree.roots[index]
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
				tagToName[tag] = names + "、" + name
			} else {
				tagToName[tag] = name
			}
		}
	}

	return tagToName
}

// export time and location
type UrlConfig struct {
	Loc  string
	Date time.Time
}

func (urlConfig *UrlConfig) ToDailyUrl() string {
	return "/log?date=" + urlConfig.Date.Format("2006-01-02") + "&loc=" + urlConfig.Loc
}

func (urlConfig *UrlConfig) ToWeekDataUrl() string {
	return "/loc/" + urlConfig.Loc + "/week/" + urlConfig.Date.Format("2006-01-02")
}

func (urlConfig *UrlConfig) ToYearStatsUrl() string {
	return "/loc/" + urlConfig.Loc + "/year/" + strconv.FormatInt(int64(urlConfig.Date.Year()), 10)
}

type RequestType int

const (
	RequestDailyLog RequestType = iota
	RequestWeekData
	RequestYearData
)

func (urlConfig *UrlConfig) Explain(reqType RequestType) string {
	switch reqType {
	case RequestDailyLog:
		return fmt.Sprintf("获取 %v 心栈 %v 的数据", bucketMap[urlConfig.Loc], urlConfig.Date.Format("2006-01-02"))
	case RequestWeekData:
		return fmt.Sprintf("获取 %v 心栈 %v 那一周的数据", bucketMap[urlConfig.Loc], urlConfig.Date.Format("2006-01-02"))
	case RequestYearData:
		return fmt.Sprintf("获取 %v 心栈 %v 年的数据", bucketMap[urlConfig.Loc], urlConfig.Date.Year())
	}

	return ""
}

type YearStats struct {
	CupSize     int `json:"cupSize"`
	NumOfTime   int `json:"numOfTime"`
	NumOfPeople int `json:"numOfPeople"`
	NumOfNew    int `json:"numOfNew"`
}

type WeekStats struct {
	YearStats
	Names []string
}

func calcWeekData(logList []Body) *WeekStats {
	nameMap := map[string]bool{}
	weekStats := &WeekStats{
		YearStats{0, 0, 0, 0},
		nil,
	}
	for _, log := range logList {
		weekStats.CupSize += log.CupSize
		weekStats.NumOfTime += len(log.Names)
		for _, tags := range log.Tags {
			if strings.Contains(tags, "新人") {
				weekStats.NumOfNew += 1
			}
		}
		for _, name := range log.Names {
			nameMap[name] = true
		}
	}
	weekStats.NumOfPeople = len(nameMap)
	weekStats.Names = make([]string, 0, len(nameMap))
	for name := range nameMap {
		weekStats.Names = append(weekStats.Names, name)
	}

	return weekStats
}
