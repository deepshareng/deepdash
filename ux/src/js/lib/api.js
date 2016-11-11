function getUsername(success, error) {
    $.ajax({
        type: "GET",
        url: "/session/getuser",
        dataType: "json",
        success: success,
        error: error,
    });
}

function getUserApps(success, error) {
    $.ajax({
        type: "GET",
        url: "/apps",
        dataType: "json",
        success: success,
        error: error,
    });
}

function setSessionAppId(appid, success, error) {
    $.ajax({
        type: "POST",
        url: "/session/setappid/" + appid,
        dataType: "json",
        data: {myappid: appid},
        success: success,
        error: error,
    });
}

function getAbstract(appid, success, error) {
    $.ajax({
        type: "GET",
        // FIXME: the tedius endpoint
        url: "/apps/" + appid + "/statistics",
        data: {
            channel: 'all',
            gran: 't',
            limit: 1,
            event: 'basic-event,url-tem-overall-value,sharelink-tem-overall-value',
        },
        dataType: "json",
        success: success,
        error: error,
    });
}

function getDist(appid, gran, success, error) {
    // gran: [d|w|m]
    $.ajax({
        type: "GET",
        url: "/apps/" + appid + "/statistics",
        data: {
            channel: 'all',
            gran: gran,
            limit: 10,
            groupby: 'os',
            event: 'basic-event,url-tem-overall-value,sharelink-tem-overall-value',
        },
        dataType: "json",
        success: success,
        error: error,
    });
}


function getChannelStatistics(appid, data, before, success, error, complete) {
    $.ajax({
        type: "GET",
        url: "/apps/" + appid + "/channels/statistics",
        data: data,
        dataType: "json",
        beforeSend: before,
        success: success,
        error: error,
        complete: complete,
    });
    /*
    data: {
        event: 'match/install_with_params,match/open_with_params,3-day-retention,7-day-retention', 
    },
    */
}

function getEvents(appid, success, error) {
    $.ajax({
        type: "GET",
        url: "/apps/" + appid + "/events",
        dataType: "json",
        success: success,
        error: error,
    });
}

function getChannelList(appid, success, error) {
    $.ajax({
        type: "GET",
        url: "/apps/" + appid + "/channels",
        dataType: "json",
        success: success,
        error: error,
    });
}

function getSmsList(appid, success, error) {
    $.ajax({
        type: "GET",
        url: "/apps/" + appid + "/smses",
        dataType: "json",
        success: success,
        error: error,
    });
}

function getAppInfo(appid, success, error) {
    $.ajax({
        type: "GET",
        url: "/apps/" + appid,
        dataType: "json",
        success: success,
        error: error,
    });
}

function generateChannelUrl(appid, param, success, error, complete) {
    $.ajax({
        type: "POST",
        url: "https://fds.so/v2/url/" + appid,
        dataType: "json",
        data: JSON.stringify(param),
        success: success,
        error: error,
        complete: complete,
    });
}

function insertChannel(appid, param, success, error, complete) {
    $.ajax({
        type: "POST",
        url: "/apps/" + appid + "/channels",
        data: param,
        success: success,
        error: error,
        complete: complete,
    });
}

function deleteChannel(appid, channelName, success, error, complete) {
    $.ajax({
        type: "DELETE",
        url: "/apps/" + appid + "/channels/" + channelName,
        success: success,
        error: error,
        complete: complete,
    });
}

function insertSms(appid, param, success, error, complete) {
    $.ajax({
        type: "POST",
        url: "/apps/" + appid + "/smses",
        data: param,
        success: success,
        error: error,
        complete: complete,
    });
}

function deleteSms(appid, smsid, success, error, complete) {
    $.ajax({
        type: "DELETE",
        url: "/apps/" + appid + "/smses/" + smsid,
        success: success,
        error: error,
        complete: complete,
    });
}

function updateSms(appid, smsid, content, success, error, complete) {
    $.ajax({
        type: "PUT",
        url: "/apps/" + appid + "/smses/" + smsid,
        data: {'content': content},
        success: success,
        error: error,
        complete: complete,
    });
}
function addApp(appname, success, error) {
    $.ajax({
        type: "POST",
        url: "/appid",
        dataType: "json",
        data: {name: appname},
        success: success,
        error: error,
    });
}

function createApp(param, success, error, before, complete) {
    $.ajax({
        type: "POST",
        url: "/apps",
        data: param,
        beforeSend: before,
        success: success,
        error: error,
        complete: complete,
    });
}

function updateAppInfo(appid, param, success, error) {
    $.ajax({
        type: "PUT",
        url: "/apps/" + appid,
        data: param,
        success: success,
        error: error,
    });
}

function deleteApp(appid, success, error) {
    $.ajax({
        type: "DELETE",
        url: "/apps/" + appid,
        success: success,
        error: error,
    });
}

function saveCallbackUrl(appid, param, success, error) {
    $.ajax({
        type: "PUT",
        url: "/apps/" + appid + "/url",
        data: param,
        success: success,
        error: error,
    });
}

function testCallback(appid, param, success, error) {
    $.ajax({
        type: "POST",
        url: "/apps/" + appid + "/callback",
        dataType: "json",
        data: param,
        success: success,
        error: error,
    });

}

function uploadImage(appid, formData, success, error, complete) {
    $.ajax({
        type: 'POST',
        url: '/apps/' + appid + '/uploadimage',
        dataType: 'json',
        data: formData,
        // !! below two
        processData: false,
        contentType: false,
        success: success,
        error: error,
        complete: complete,
    });
}

function login(username, password, success, error, complete) {
    $.ajax({
        type: "POST",
        url: "/authorization",
        data: {
            "username": username,
            "password": password,
        },
        success: success,
        error: error,
        complete: complete,
    });
}


var DS_CONST_API_CODE_PARAM_INVALID = 6000;
var DS_CONST_API_CODE_AUTH_FAIL = 1101;


module.exports = {
    _version_: '1.0',

    gotoLoginPage: function() {
        gotoUrl('./login.html');
    },

    getUsername: getUsername,
    getUserApps: getUserApps,
    getAbstract: getAbstract,
    getDist: getDist,

    setSessionAppId: setSessionAppId,

    getAppInfo: getAppInfo,
    getEvents: getEvents,
    getChannelList: getChannelList,
    getChannelStatistics: getChannelStatistics,
    insertChannel: insertChannel,
    deleteChannel: deleteChannel,
    getSmsList: getSmsList,
    insertSms: insertSms,
    deleteSms: deleteSms,
    updateSms: updateSms,
    generateChannelUrl: generateChannelUrl,

    addApp: addApp,
    createApp: createApp,
    updateAppInfo: updateAppInfo,
    deleteApp: deleteApp,
    saveCallbackUrl: saveCallbackUrl,
    testCallback: testCallback,

    uploadImage: uploadImage,
    login: login,

    DS_CONST_API_CODE_PARAM_INVALID: DS_CONST_API_CODE_PARAM_INVALID,
    DS_CONST_API_CODE_AUTH_FAIL: DS_CONST_API_CODE_AUTH_FAIL,
};
