package main

import (
	"fmt"
	"github.com/lxn/walk"
	"strconv"
	"strings"
	"time"
)

type LocPair struct {
	key, value string
}

func (loc *LocPair) ToString() string {
	return "NO." + loc.key + " " + loc.value
}

var bucketSlice = make([]LocPair, 0, len(bucketMap))

func locIndexOf(arr []LocPair, key string) int {
	for index, loc := range arr {
		if loc.key == key {
			return index
		}
	}
	return 0
}

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

type LocListBoxAdapter struct {
	view  *walk.ListBox
	model *LocListModel
}

type LocComboBoxAdapter struct {
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
type YearItem struct {
	year     int
	children []*MonthItem
}

func (tree *YearItem) Text() string {
	return strconv.FormatInt(int64(tree.year), 10) + "年"
}

func (*YearItem) Parent() walk.TreeItem {
	return nil
}

func (tree *YearItem) ChildCount() int {
	return len(tree.children)
}

func (tree *YearItem) ChildAt(index int) walk.TreeItem {
	return tree.children[index]
}

type MonthItem struct {
	year     int
	month    int
	parent   *YearItem
	children []*DayItem
}

func (tree *MonthItem) Text() string {
	return fmt.Sprintf("%v年%v月", tree.year, tree.month)
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
	t      time.Time
	parent *MonthItem
}

func (tree *DayItem) Text() string {
	return tree.t.Format("2006-01-02")
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
	roots []*YearItem
}

func newEmptyTreeModel() *LocTreeModel {
	roots := make([]*YearItem, 0)

	return &LocTreeModel{roots: roots}
}

func addToRoot(roots []*YearItem, t time.Time) []*YearItem {
	year := t.Year()
	month := int(t.Month())

	var yearItem *YearItem
	for _, yearChild := range roots {
		if yearChild.year == year {
			yearItem = yearChild
			break
		}
	}
	if yearItem == nil {
		yearItem = &YearItem{year, make([]*MonthItem, 0)}
		roots = append(roots, yearItem)
	}

	var monthItem *MonthItem
	for _, monthChild := range yearItem.children {
		if monthChild.month == month {
			monthItem = monthChild
			break
		}
	}
	if monthItem == nil {
		monthItem = &MonthItem{year, month, yearItem, make([]*DayItem, 0)}
		yearItem.children = append(yearItem.children, monthItem)
	}

	dayItem := &DayItem{t, monthItem}
	monthItem.children = append(monthItem.children, dayItem)

	return roots
}

func treeModelFromList(keys []string) *LocTreeModel {
	timeList := make([]time.Time, 0, len(keys))
	for _, key := range keys {
		t, err := time.Parse("2006-01-02", key)
		if err != nil {
			continue
		} else {
			timeList = append(timeList, t)
		}
	}

	if len(timeList) == 0 {
		return newEmptyTreeModel()
	}

	roots := make([]*YearItem, 0)
	for _, t := range timeList {
		roots = addToRoot(roots, t)
	}

	return &LocTreeModel{roots: roots}
}

func (*LocTreeModel) LazyPopulation() bool {
	return false
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

func weekToCN(w time.Weekday) string {
	switch w {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	default:
		return "周日"
	}
}

func (log *Body) getPreview(urlConfig *UrlConfig) string {
	t := urlConfig.Date
	loc := bucketMap[urlConfig.Loc]
	tagToName := log.remixTagTable()
	ret := fmt.Sprintf("标题：仁爱 %v心栈 %v年%v月%v日奉粥日志+题目\r\n", loc, t.Year(), int(t.Month()), t.Day())
	ret += fmt.Sprintf("奉粥日期：%v年%v月%v日（%v）\r\n", t.Year(), int(t.Month()), t.Day(), weekToCN(t.Weekday()))
	ret += fmt.Sprintf("日负责人：%v\r\n", strings.Join(tagToName["负责人"], "、"))
	ret += fmt.Sprintf("签到：%v\r\n", strings.Join(tagToName["签到"], "、"))
	ret += fmt.Sprintf("熬粥：%v\r\n", strings.Join(tagToName["熬粥"], "、"))
	ret += fmt.Sprintf("前行：%v\r\n", strings.Join(tagToName["前行"], "、"))
	ret += fmt.Sprintf("杯数：%v 杯\r\n", log.CupSize)
	ret += fmt.Sprintf("人数：%v 人\r\n", len(log.Names))
	ret += fmt.Sprintf("新人数：%v 人，%v\r\n", len(tagToName["新人"]), strings.Join(tagToName["新人"], "、"))
	ret += fmt.Sprintf("摄影：%v\r\n", strings.Join(tagToName["摄影"], "、"))
	ret += fmt.Sprintf("日志：无\r\n")
	ret += fmt.Sprintf("文宣：\r\n")
	ret += fmt.Sprintf("结行：%v\r\n", strings.Join(tagToName["结行"], "、"))
	ret += fmt.Sprintf("后勤：%v\r\n", strings.Join(tagToName["后勤"], "、"))
	ret += fmt.Sprintf("环保：%v\r\n", strings.Join(tagToName["环保"], "、"))
	ret += fmt.Sprintf("奉粥：%v\r\n", strings.Join(tagToName["奉粥"], "、"))

	return ret
}

func (log *Body) remixTagTable() map[string][]string {
	tagToName := map[string][]string{
		"负责人": {},
		"签到":  {},
		"熬粥":  {},
		"前行":  {},
		"环保":  {},
		"摄影":  {},
		"新人":  {},
		"结行":  {},
		"后勤":  {},
		"奉粥":  {},
	}
	for index, tags := range log.Tags {
		name := log.Names[index]
		// 所有没有参加环保的都是奉粥人员
		ok := strings.Contains(tags, "环保")
		if !ok {
			tagToName["奉粥"] = append(tagToName["奉粥"], name)
		}
		tagList := strings.Split(tags, "|")
		for _, tag := range tagList {
			nameList, exist := tagToName[tag]
			if exist {
				tagToName[tag] = append(nameList, name)
			}
		}
	}

	return tagToName
}

func (*Body) formatTags(tag string) string {
	if tag == "" {
		return tag
	}

	return "（" + strings.Replace(tag, "|", "、", -1) + "）"
}

func (log *Body) getNamesWithTags() []string {
	nameList := make([]string, 0, len(log.Names))

	for index, name := range log.Names {
		nameList = append(nameList, name+log.formatTags(log.Tags[index]))
	}

	return nameList
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
	CupSize   int      `json:"cupSize"`
	NumOfTime int      `json:"numOfTime"`
	NumOfNew  int      `json:"numOfNew"`
	Names     []string `json:"names"`
}

type WeekStats struct {
	YearStats
	NumOfPeople int `json:"numOfPeople"`
}

func calcWeekData(logList []Body) *WeekStats {
	nameMap := map[string]bool{}
	weekStats := &WeekStats{
		YearStats{0, 0, 0, nil},
		0,
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
