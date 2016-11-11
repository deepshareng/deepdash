package client

import "github.com/MISingularity/deepdash/storage"

func GetUniqueID(githubname, token string) (string, error) {
	err := storage.MongoUS.UpdateUser(storage.UserFormat{Githubname: githubname}, storage.UserFormat{Githubname: githubname, Token: token}, true)
	if err != nil {
		return "", err
	}
	results, err := storage.MongoUS.GetUser(storage.UserFormat{Githubname: githubname}, false, []string{}, []string{})
	if err != nil {
		return "", err
	}
	res := results[0]
	return res.Username, nil
}

func CheckUserExist(username string) (bool, error) {
	matchExistUser := storage.UserFormat{
		Username: username,
	}
	existUser, err := storage.MongoUS.GetUser(matchExistUser, false, []string{}, []string{})
	if err != nil {
		return false, nil
	}
	if len(existUser) > 0 {
		if existUser[0].Activate == "1" {
			return true, nil
		}
		matchExistToken := storage.TokenFormat{
			AccountID: username,
		}
		existToken, err := storage.MongoTS.GetToken(matchExistToken, false)
		if len(existToken) > 0 {
			return true, nil
		}
		return false, err
	}
	return false, err
}

// check if username is not active
// only when user exist and active is not equle to "1" return true
// if user exist return user
// if password is right return 'true'
func CheckUserIsNotActive(username string, password string) (bool, storage.UserFormat, bool, error) {
	matchExistUser := storage.UserFormat{
		Username: username,
	}
	rightPassword := false
	existUser, err := storage.MongoUS.GetUser(matchExistUser, false, []string{}, []string{})
	if err != nil {
		return false, storage.UserFormat{}, rightPassword, err
	}
	if len(existUser) > 0 {
		if password == existUser[0].Password {
			rightPassword = true
		}
		if existUser[0].Activate != "1" {
			return true, existUser[0], rightPassword, nil
		} else {
			return false, existUser[0], rightPassword, nil
		}
	}
	return false, storage.UserFormat{}, rightPassword, nil
}
