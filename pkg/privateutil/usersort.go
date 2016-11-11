package privateutil

import (
	"sort"
	"strconv"

	"github.com/MISingularity/deepdash/api"
)

type UserStatusList []api.PrivateUserStatus

func (a UserStatusList) Len() int      { return len(a) }
func (a UserStatusList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a UserStatusList) Less(i, j int) bool {
	if a[i].CreateAt != "" && a[j].CreateAt != "" {
		iCreateAt, _ := strconv.ParseInt(a[i].CreateAt, 10, 64)
		jCreateAt, _ := strconv.ParseInt(a[j].CreateAt, 10, 64)
		return iCreateAt >= jCreateAt
	}
	if a[i].CreateAt == "" && a[j].CreateAt == "" {
		return a[i].Username < a[j].Username
	}
	if a[i].CreateAt != "" {
		return true
	}
	return false
}

func SortPrivateUserStatus(list UserStatusList) UserStatusList {
	sort.Sort(list)
	return list
}
