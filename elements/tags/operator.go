package tags

import (
	"sort"
)

// TODO!!!
type Item = string

//type Tags []Item
//func (tags Tags) String() string {
//	//var tagsStr []string
//	//for _, tag := range tags {
//	//
//	//}
//}

type StatMap map[Item]int64
type Stat struct {
	Tag   Item
	Count int64
}
type Stats []Stat

func (ts StatMap) List(sortByCount bool) Stats {
	var tagsStatList Stats
	for t, c := range ts {
		tagsStatList = append(tagsStatList, Stat{t, c})
	}

	if sortByCount {
		sort.Slice(tagsStatList, func(i, j int) bool { return tagsStatList[i].Count >= tagsStatList[j].Count })
	} else {
		sort.Slice(tagsStatList, func(i, j int) bool { return tagsStatList[i].Tag <= tagsStatList[j].Tag })
	}

	return tagsStatList
}
