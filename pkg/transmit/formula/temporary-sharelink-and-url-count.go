package formula

import "strconv"

type SharelinkTemFormula struct{}

func (f *SharelinkTemFormula) Requisite(param string, gran string) []string {
	return []string{"sharelink:", "sharelink:/d/_no_ds_tag", "/v2/dsactions/ds_jssdk_click"}
}

func (f *SharelinkTemFormula) Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string {
	res := make(map[string][]string)
	res["sharelink-tem-overall-value"] = make([]string, limit)
	for i := 0; i < limit; i++ {
		d := 0
		if i < len(data["sharelink:"]) {
			d += data["sharelink:"][i]
		}
		if i < len(data["sharelink:/d/_no_ds_tag"]) {
			d += data["sharelink:/d/_no_ds_tag"][i]
		}
		if i < len(data["/v2/dsactions/ds_jssdk_click"]) {
			d += data["/v2/dsactions/ds_jssdk_click"][i]
		}
		res["sharelink-tem-overall-value"][i] = strconv.Itoa(d)
	}
	return res
}

type GenURLTemFormula struct{}

func (f *GenURLTemFormula) Requisite(param string, gran string) []string {
	return []string{"v2/url-old", "/v2/url/", "sharelink:/v2/jsapi/"}
}

func (f *GenURLTemFormula) Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string {
	res := make(map[string][]string)
	res["url-tem-overall-value"] = make([]string, limit)
	for i := 0; i < limit; i++ {
		d := 0
		if i < len(data["v2/url-old"]) {
			d += data["v2/url-old"][i]
		}
		if i < len(data["/v2/url/"]) {
			d += data["/v2/url/"][i]
		}
		if i < len(data["sharelink:/v2/jsapi/"]) {
			d += data["sharelink:/v2/jsapi/"][i]
		}
		res["url-tem-overall-value"][i] = strconv.Itoa(d)
	}
	return res
}
