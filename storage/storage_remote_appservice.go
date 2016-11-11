package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

type RemoteAppInfoService struct {
	appinfourl string
}

func NewRemoteAppInfoService(addr string) AppStorageService {
	return &RemoteAppInfoService{addr}
}

func (m *RemoteAppInfoService) GetApplist(match AppFormat) ([]AppFormat, error) {
	return []AppFormat{}, nil
}

func (m *RemoteAppInfoService) GetApplistBson(match bson.M) ([]AppFormat, error) {
	// Not implemented yet !
	return []AppFormat{}, nil
}

func (m *RemoteAppInfoService) GetApp(match AppFormat) (AppFormat, error) {
	return AppFormat{}, nil
}

func (m *RemoteAppInfoService) GetAppBson(match bson.M) (AppFormat, error) {
	// Not implemented yet !
	return AppFormat{}, nil
}

func (m *RemoteAppInfoService) InsertApp(insert AppFormat) error {
	androidInfo := AppAndroidInfo{}
	iosInfo := AppIosInfo{}
	userConf := UserConfig{}

	androidInfo.DownloadUrl = insert.AndroidDownloadUrl
	androidInfo.Host = insert.AndroidHost
	androidInfo.Pkg = insert.AndroidPkgname
	androidInfo.Scheme = insert.AndroidScheme
	androidInfo.IsDownloadDirectly = insert.AndroidIsDownloadDirectly == "true"
	androidInfo.YYBEnable = insert.AndroidYYBEnable == "true"

	iosInfo.DownloadUrl = insert.IosDownloadUrl
	iosInfo.Scheme = insert.IosScheme
	iosInfo.UniversalLinkEnable = insert.IosUniversalLinkEnable == "true"
	iosInfo.BundleID = insert.IosBundler
	iosInfo.TeamID = insert.IosTeamID
	iosInfo.YYBEnableBelow9 = insert.IosYYBEnableBelow9 == "true"
	iosInfo.YYBEnableAbove9 = insert.IosYYBEnableAbove9 == "true"
	iosInfo.ForceDownload = insert.ForceDownload == "true"

	userConf.BgWeChatAndroidTipUrl = insert.UserConfBgWeChatAndroidTipUrl
	userConf.BgWeChatIosTipUrl = insert.UserConfBgWeChatIosTipUrl

	info := AppInfo{
		AppID:     insert.AppID,
		ShortID:   insert.ShortID,
		AppName:   insert.AppName,
		Android:   androidInfo,
		Ios:       iosInfo,
		YYBUrl:    insert.YYBurl,
		YYBEnable: insert.AndroidYYBEnable == "true",
		IconUrl:   insert.IconUrl,
		Theme:     insert.Theme,
		UserConf:  userConf,
	}
	return m.appinfoPUTRequest(&insert, &info)
}

func (m *RemoteAppInfoService) UpdateApp(match AppFormat, update AppFormat, upsert bool) error {
	return m.InsertApp(update)
}

func (m *RemoteAppInfoService) UpdateAppBson(match bson.M, update bson.M, upsert bool) error {
	// Not implemented
	return nil
}

func (m *RemoteAppInfoService) DeleteApp(appid string) error {
	// Not implemented
	return nil
}

func (m *RemoteAppInfoService) appinfoPUTRequest(putApp *AppFormat, putAppinfo *AppInfo) error {
	jsonStr, err := json.Marshal(putAppinfo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", m.appinfourl+"/v2/appinfo/"+putApp.AppID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Remote appinfo server submit failed")
	}

	defer resp.Body.Close()
	return err
}

func (m *RemoteAppInfoService) GetAppleAppSiteAssociationInfo() ([]map[string]string, error) {
	return []map[string]string{}, nil
}
