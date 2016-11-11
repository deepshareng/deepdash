package cert

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/storage"
)

var iOSLock, androidLock sync.Mutex

// generate json file for universal link
func GenerateCertFile(appId string, teamId string, bundlerId string) {
	// mutex to lock
	iOSLock.Lock()
	defer iOSLock.Unlock()

	curDir := getcurdir()
	/*
		file, err := ioutil.ReadFile(curDir + "/../app_site/apple-app-site-association-unsigned")
		if err != nil {
			log.Errorf("[Cert Resource] Find file <apple-app-site-association-unsigned> failed! Err Msg=%s", err)
			return
		}

		var universallinkJson UniversalLinkJson
		json.Unmarshal(file, &universallinkJson)

		// add new element
		comboId := teamId + "." + bundlerId
		path := "/d/" + appId + "/*"
		s := Detail{AppID: comboId, Paths: []string{path}}
		universallinkJson.Applinks.Details = append(universallinkJson.Applinks.Details, s)

		newResult, err := json.Marshal(universallinkJson)
		if err != nil {
			log.Errorf("[Cert Resource] Marshal upsert info failed! Err Msg=%v", err)
			return
		}
	*/

	// FIXME: when two many apps
	result, err := storage.MongoAS.GetAppleAppSiteAssociationInfo()
	/*
		httperror := errorutil.ProcessMongoError(f, err)
		if httperror {
			return
		}*/
	r, err := json.Marshal(result)
	log.Debugf("db result: %s\n", r)
	// Merge mutiple paths of the same app(same teamid & bundleid)
	rawlist := map[string]interface{}{}

	applist := []map[string]interface{}{}

	for i, _ := range result {
		key := result[i]["iosteamid"] + "." + result[i]["iosbundler"]
		fullPath := "/d/" + result[i]["appid"] + "/*"
		miniPath := "/d/" + result[i]["shortid"] + "/*"

		if _, ok := rawlist[key]; !ok {
			rawlist[key] = []string{}
		}

		if result[i]["shortid"] != "" {
			rawlist[key] = append(rawlist[key].([]string), miniPath)
		}
		rawlist[key] = append(rawlist[key].([]string), fullPath)
	}

	for k, v := range rawlist {
		app := map[string]interface{}{}
		app["appID"] = k
		app["paths"] = v
		applist = append(applist, app)
	}
	ujson := map[string]interface{}{
		"applinks": map[string]interface{}{
			"apps":    []string{},
			"details": applist,
		},
	}
	newResult, err := json.Marshal(ujson)
	log.Infof("The generated apple-app-site-association file: %s\n", newResult)
	err = ioutil.WriteFile(curDir+"/../app_site/apple-app-site-association-unsigned", newResult, 0666)
	if err != nil {
		log.Errorf("[Cert Resource] Failed to generate new json file, error: %v\n", err)
		return
	}

	// write with plain json file
	err = ioutil.WriteFile(curDir+"/../app_site/apple-app-site-association", newResult, 0666)
	if err != nil {
		log.Errorf("[Cert Resource] Failed to generate apple-app-site-association file, error: %v\n", err)
		return
	}
	/*
		* No need to sign:
		* https://developer.apple.com/library/ios/documentation/General/Conceptual/AppSearch/UniversalLinks.html
		arv := []string{"smime", "-sign", "-nodetach", "-in", curDir + "/../app_site/apple-app-site-association-unsigned", "-out", curDir + "/../app_site/apple-app-site-association", "-outform", "DER", "-inkey", curDir + "/fds.so.key", "-signer", curDir + "/fds.so.crt"}
		cmd := exec.Command("openssl", arv...)

		if err := cmd.Run(); err != nil {
			log.Errorf("[Cert Resource] Failed to sign the json file! Err Msg=%v", err)
		} else {
			log.Infof("[Cert Resource] Success to sign the json file!")
		}
	*/
}

// generate assetlinks json file for app link
func GenerateAppLinkJsonFile(packageName string, sha256 string) {
	log.Infof("[Cert Resource] Start to generate new app link json file!")
	// mutex to lock
	androidLock.Lock()
	defer androidLock.Unlock()

	curDir := getcurdir()
	file, err := ioutil.ReadFile(curDir + "/../app_site/.well-known/assetlinks_full.json")
	if err != nil {
		log.Errorf("[Cert Resource] Find file <assetlinks> failed! Err Msg=%s", err)
		return
	}

	var appLinkJson AppLinkJson
	json.Unmarshal(file, &appLinkJson)
	log.Infof("[Cert Resource] AppLinkJson File Content : %v", appLinkJson)
	// add new element

	target := Target{Namespace: "android_app", PackageName: packageName, SHA256: []string{sha256}}
	applink := AppLink{Relation: []string{"delegate_permission/common.handle_all_urls"}, Target: target}

	appLinkJson.AppLink = append(appLinkJson.AppLink, applink)

	newResult, err := json.Marshal(appLinkJson)
	if err != nil {
		log.Errorf("[Cert Resource] Marshal upsert info failed! Err Msg=%v", err)
		return
	}

	err = ioutil.WriteFile(curDir+"/../app_site/.well-known/assetlinks_full.json", newResult, 0666)
	if err != nil {
		log.Errorf("[Cert Resource] Failed to generate new json file, error: %v\n", err)
		return
	}

	//
	newResult, err = json.Marshal(appLinkJson.AppLink)
	if err != nil {
		log.Errorf("[Cert Resource] Marshal upsert info failed! Err Msg=%v", err)
		return
	}

	err = ioutil.WriteFile(curDir+"/../app_site/.well-known/assetlinks.json", newResult, 0666)
	if err != nil {
		log.Errorf("[Cert Resource] Failed to generate new json file, error: %v\n", err)
		return
	}
}
func getcurdir() string {
	//get absolute path, which will be used to locate html and js files
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		panic(err)
	}
	return dir
}
