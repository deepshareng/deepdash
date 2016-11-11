package storage

import "time"

func PrepareTokeItems(ts TokenStorageService) {
	ts.InsertToken(TokenFormat{
		PasswordID: "PasswordID1",
		Token:      "token1",
		CreateAt:   time.Date(2015, time.March, 1, 1, 2, 3, 0, time.Local),
	})
	ts.InsertToken(TokenFormat{
		PasswordID: "PasswordID2",
		Token:      "token2",
		CreateAt:   time.Date(2015, time.March, 2, 1, 2, 3, 0, time.Local),
	})
	ts.InsertToken(TokenFormat{
		PasswordID: "PasswordID3",
		Token:      "token3",
		CreateAt:   time.Date(2015, time.March, 3, 1, 2, 3, 0, time.Local),
	})
}

func PrepareChannelItems(us ChannelStorageService) {
	us.InsertChannel(AppChannel{
		AppID:       "appid1",
		Channelname: "channelname1",
	})

	us.InsertChannel(AppChannel{
		AppID:       "appid2",
		Channelname: "channelname2",
	})

	us.InsertChannel(AppChannel{
		AppID:       "appid3",
		Channelname: "channelname3",
	})

	us.InsertChannel(AppChannel{
		AppID:       "appid4",
		Channelname: "channelname4",
	})
}

func PrepareAppItems(us AppStorageService) {
	us.InsertApp(AppFormat{
		AppID:     "appid1",
		AccountID: "account1",
		AppName:   "app1",
		PkgName:   "pkg1",
	})

	us.InsertApp(AppFormat{
		AppID:     "appid2",
		AccountID: "account1",
		AppName:   "app2",
		PkgName:   "pkg2",
	})

	us.InsertApp(AppFormat{
		AppID:     "appid3",
		AccountID: "account3",
		AppName:   "app3",
		PkgName:   "pkg3",
	})

	us.InsertApp(AppFormat{
		AppID:     "appid4",
		AccountID: "account4",
		AppName:   "app4",
		PkgName:   "pkg4",
	})
}

func PrepareUserItems(us UserStorageService) {
	us.InsertUser(UserFormat{
		Username:   "test1",
		Password:   "123",
		Githubname: "git1",
		Activate:   "1",
	})
	us.InsertUser(UserFormat{
		Username:   "test2",
		Password:   "123",
		Githubname: "git2",
		Activate:   "1",
	})
	us.InsertUser(UserFormat{
		Username:   "test3",
		Password:   "123",
		Githubname: "git3",
		Activate:   "1",
	})
	us.InsertUser(UserFormat{
		Username:   "test4",
		Password:   "123",
		Githubname: "git4",
	})
}

func PrepareAppSelectedItems(as AppSelectedItemsService) {
	as.InsertAppSelectedItems(AppSelectedItems{
		AppID:     "appid1",
		AccountID: "acc1",
		Events:    []string{"a1", "b1", "c1"},
		Displays:  []string{"da1", "db1", "dc1"},
	})

	as.InsertAppSelectedItems(AppSelectedItems{
		AppID:     "appid2",
		AccountID: "acc2",
		Events:    []string{"a2", "b2", "c2"},
		Displays:  []string{"da2", "db2", "dc2"},
	})

	as.InsertAppSelectedItems(AppSelectedItems{
		AppID:     "appid3",
		AccountID: "acc3",
		Events:    []string{"a3", "b3", "c3"},
		Displays:  []string{"da3", "db3", "dc3"},
	})
}
