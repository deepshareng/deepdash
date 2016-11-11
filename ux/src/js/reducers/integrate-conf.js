var ConfReducer = function(state={}, action) {
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {
            appid: sessionAppId,
            appName: '',
            yyburl: '',
            theme: '0',
            downloadTitle: '',
            downloadMsg: '',
            iconUrl: '',

            userConfBgWeChatAndroidTipUrl: '',
            userConfBgWeChatIosTipUrl: '',
            iosBundler: '',
            iosScheme: 'ds' + sessionAppId,
            iosDownloadUrl: '',
            iosTeamID: '',

            iosYYBEnableAbove9: 'false',
            iosYYBEnableBelow9: 'false',
            forceDownload: 'false',

            androidPkgname: '',
            androidHost: '',
            androidYYBEnable: 'false',

            androidIsDownloadDirectly: '', 
            androidDownloadUrl: '',
            androidSHA256: '',
            androidScheme: 'ds' + sessionAppId,
        };
    } else {
        if (action.type === 'UPDATE_APP_CONF') {
            return Object.assign({}, state, action.conf);
        } else {
            return state;
        }
    }
};

module.exports = ConfReducer;
