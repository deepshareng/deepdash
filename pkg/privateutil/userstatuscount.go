package privateutil

import "github.com/MISingularity/deepdash/api"

func GetStatusCount(userStatus UserStatusList) (lastWeekCount api.StatusCounter, accountCount api.StatusCounter) {
	for _, v := range userStatus {
		switch v.AccountStatus {
		case "集成完毕":
			accountCount.Integrated++
		case "集成中":
			accountCount.Integrating++
		case "集成失败":
			accountCount.IntegrateFailed++
		case "注册完毕":
			accountCount.Registered++
		case "冻结中":
			accountCount.Freeze++
		}
		switch v.LastWeekAccountStatus {
		case "集成完毕":
			lastWeekCount.Integrated++
		case "集成中":
			lastWeekCount.Integrating++
		case "集成失败":
			lastWeekCount.IntegrateFailed++
		case "注册完毕":
			lastWeekCount.Registered++
		case "冻结中":
			lastWeekCount.Freeze++
		}
	}
	return lastWeekCount, accountCount
}

func GetAppStatusCount(appStatus []api.PrivateAppStatus) (accountCount api.StatusCounter) {
	for _, v := range appStatus {
		switch v.AccountStatus {
		case "集成完毕":
			accountCount.Integrated++
		case "集成中":
			accountCount.Integrating++
		case "集成失败":
			accountCount.IntegrateFailed++
		case "注册完毕":
			accountCount.Registered++
		case "冻结中":
			accountCount.Freeze++
		}
	}
	return accountCount
}
