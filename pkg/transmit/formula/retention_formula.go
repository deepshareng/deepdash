package formula

import (
	"strconv"

	"github.com/MISingularity/deepshare2/api"
)

type TotalRetentionFormula struct{}

func paramExist(data map[string][]int, param string) bool {
	_, ok := data[param]
	return ok
}

func (f *TotalRetentionFormula) Requisite(param string, gran string) []string {
	return []string{api.RetentionAmountPrefix + param + "_retention-rate", api.RetentionAmountPrefix + param + "_retention-day"}
}

func (f *TotalRetentionFormula) Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string {
	res := make(map[string][]string)
	numerator := api.RetentionAmountPrefix + param + "_retention-rate"
	denominator := api.RetentionAmountPrefix + param + "_retention-day"
	for i := 0; i < limit; i++ {
		res[param] = make([]string, limit)
		if paramExist(data, denominator) && data[denominator][i] > 0 {
			res[param][i] = strconv.FormatFloat(float64(data[numerator][i])/100/float64(data[denominator][i]), 'f', 2, 64)
		} else {
			res[param][i] = "0.00"
		}
	}
	return res
}

type DayRetentionFormula struct{}

func (f *DayRetentionFormula) Requisite(param string, gran string) []string {
	return []string{api.RetentionPrefix + param, api.RetentionPrefix + param + "_install"}
}

func (f *DayRetentionFormula) Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string {
	res := make(map[string][]string)
	res[param] = make([]string, limit)
	for i := 0; i < limit; i++ {
		if paramExist(data, param+"_install") && data[param+"_install"][i] > 0 {
			res[param][i] = strconv.FormatFloat(float64(data[param][i])*100/float64(data[param+"_install"][i]), 'f', 2, 64)
		} else {
			res[param][i] = "0.00"
		}
	}
	return res
}

type RetentionFormula struct{}

func (f *RetentionFormula) Requisite(param string, gran string) []string {
	if gran == "d" || gran == "w" {
		tem := &DayRetentionFormula{}
		return tem.Requisite(param, gran)
	}
	if gran == "t" {
		tem := &TotalRetentionFormula{}
		return tem.Requisite(param, gran)
	}
	return []string{}
}
func (f *RetentionFormula) Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string {
	if gran == "d" || gran == "w" {
		tem := &DayRetentionFormula{}
		return tem.Calculate(param, data, gran, limit)
	}
	if gran == "t" {
		tem := &TotalRetentionFormula{}
		return tem.Calculate(param, data, gran, limit)
	}
	return map[string][]string{}
}
