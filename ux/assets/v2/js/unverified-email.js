(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);throw new Error("Cannot find module '"+o+"'")}var f=n[o]={exports:{}};t[o][0].call(f.exports,function(e){var n=t[o][1][e];return s(n?n:e)},f,f.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
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

},{}],2:[function(require,module,exports){
var DSAPI = require('./api.js');
var DSUTIL = require('../lib/util.js');

function setSessionAppId() {
    DSAPI.setSessionAppId(sessionAppId);
}

function getSelectAppIndexFromSession(apps, appid) {
    for (var i = 0; i < apps.length; ++i) {
        if (apps[i].appid === appid) {
            return i;
        }
    }
    return 0;
}

var initNav = function() {
    $("[data-toggle='popover']").popover();
    $('#side-menu a').on('click', function() {
        $('#side-menu a').removeClass('active');
        $(this).addClass('active'); 
    });
};

var renderUserInfo = function () {
    DSAPI.getUsername(function(data) {
        $("#username").html(DSUTIL.escapeHtml(data.username));
    });
};

var renderApps = function(apps, appIndex, callback) {
    var tplList = [];
    $('#app-name').text(DSUTIL.escapeHtml(apps[appIndex].appname));
    for (var i in apps) {
        tplList.push('<a href="javascript:;" data-index="' + i + '">' + DSUTIL.escapeHtml(apps[i].appname) + '</a>');
    }
    $('#app-list').html(tplList.join(''));

    $('#slidebar-app-icon').attr({
        "src": apps[appIndex].iconurl || DS_CONST_ICON_DEFAULT,
    });
    $('#slidebar-app-name').html(DSUTIL.escapeHtml(apps[appIndex].appname));
    $('#navbar-app-icon').attr({
        "src": apps[appIndex].iconurl || DS_CONST_ICON_DEFAULT,
    });
    
    callback && callback();
};

var submitLogin = function(e) {
    $('#ds-login-modal-submit').off('click');

    var username = $("#username_login").val();
    var password = $("#password_login").val();

    DSAPI.login(username, password, function(data){
        if ($("#remember-me")[0].checked) {
            window.localStorage.dszz = window.btoa("yes");
            window.localStorage.dsxx = window.btoa(username);
            window.localStorage.dsyy = window.btoa(password);
        } else {
            window.localStorage.dszz = window.btoa("no");
            window.localStorage.removeItem("dsxx");
            window.localStorage.removeItem("dsyy");
        }
        // TODO: 
        hideLoginModal();
    }, function() {
        alertTips('用户名或密码错误');
    }, function() {
        $('#ds-login-modal-submit').on('click', submitLogin);
    }); 

};

var initLoginModal = function() {
    // auto complete, remember me
    if (window.atob(window.localStorage.dszz || '') === 'yes') {
        $('#remember-me')[0].checked = true;
        $("#username_login").val(window.atob(window.localStorage.dsxx));
        $("#password_login").val(window.atob(window.localStorage.dsyy));
    }

    $('#remember-me').on('change', function() {
        if (!this.checked) {
            window.localStorage.clear();
        }
    });

    $('#ds-login-modal-submit').on('click', submitLogin);
};


var getApps = function(callback) {
    DSAPI.getUserApps(function(data) {
        var appIndex = 0;
        if (JSON.stringify(data.applist) === "[]") {
            location.href = '/#/create-app';
            alertTips("您尚无任何应用，请先添加应用！");
        } else {
            if (sessionAppId === null || sessionAppId.length === 0) {
                appIndex = 0;
            } else {
                appIndex = getSelectAppIndexFromSession(data.applist, sessionAppId);
            }
            
            // FIXME: only sessionAppId changed should call callback
            if (callback) {
                callback({
                    appNames: data.applist,
                    appIndex: appIndex,
                });
            }
        }
    });
};

function selectApp(selectedAppIndex, callback) {
    appIndex = selectedAppIndex;
    renderApps(appNames, appIndex, function() {
        $('#app-list a').on('click', function() {
            selectApp($(this).data('index'), callback); 
        });
    });
    sessionAppId = appNames[appIndex].appid;
    setSessionAppId();

    callback && callback();
}

function uploadImage(imgInput, showAreaSelector, urlInputSelector, bucket, callback) {
    var image = imgInput.files[0];  
    var formData = new FormData();
    var previousIconUrl = $(urlInputSelector).val();

    if (!/image.*/.test(image.type)) {
        alertTips('您可能上传的不是图片！');
        return;
    }
    if (image.size > 5 * 1000 * 1000) {
        alertTips('您上传的图片太大了，最好的2M一下！');
        return;
    }

    formData.append('bucket', bucket);
    formData.append('uploadfile', image, image.name);
    
    // Show loading
    $(showAreaSelector).attr('src', DS_CONST_LOADING_DEFAULT); 

    DSAPI.uploadImage(sessionAppId, formData, function(data) {
        $(showAreaSelector).attr('src', data.value); 
        $(urlInputSelector).val(data.value);
        callback && callback(data);
    }, function() {
        $(showAreaSelector).attr('src', previousIconUrl);
        alertTips('上传失败，请重新上传'); 
    }, function() {
        // clear input, or choose the same file would not trigger change event
        imgInput.value = ''; 
    });

}

function showModal(title, body) {
    $('#ds-msg-modal').modal('show');
    $('#ds-msg-modal .modal-title').html(title);
    $('#ds-msg-modal .modal-body').html(body);
}

function hideModal() {
    $('#ds-msg-modal').modal('hide');
}

function showLoginModal() {
    $('#ds-login-modal').modal('show');
}
function hideLoginModal() {
    $('#ds-login-modal').modal('hide');
}

function showConfirmModal(before, cancel, ensure) {
    if (before) {
        before();
    }
    $('#ds-confirm-modal .ds-cancel').on('click', function(){
        cancel && cancel(); 
        $('#ds-confirm-modal').modal('hide');
    });
    $('#ds-confirm-modal .ds-ensure').on('click', function(){
        ensure && ensure(); 
        $('#ds-confirm-modal').modal('hide');
    });
    $('#ds-confirm-modal').modal('show');
}

function showDeleteModal(before, cancel, ensure) {
    if (before) {
        before();
    }
    $('#ds-delete-modal .ds-cancel').on('click', function(){
        cancel && cancel(); 
        $('#ds-delete-modal').modal('hide');
    });
    $('#ds-delete-modal .ds-ensure').on('click', function(){
        ensure && ensure(); 
        //$('#ds-delete-modal').modal('hide');
    });
    $('#ds-delete-modal').modal('show');
}

var tipsTimeout; 
var alertTips = function(msg, options) {
    if (!msg) {
        return;
    }

    options = options || {};
    options.level = options.level || 'info';
    options.time = options.time || 1500;

    if ($('.alert-tips-div').length > 0) {
        $('.alert-tips-div .alert-tips').html(msg);
        clearTimeout(tipsTimeout);
    } else {
        $('body').append($('<div class="alert-tips-div"><span class="alert-tips ' + options.level + '">' + msg + '</span></div>').css('display', 'none'));
        $('.alert-tips-div .alert-tips').click(function() {
            $('.alert-tips-div').fadeOut();
        });
    }
    $('.alert-tips-div').fadeIn();
    tipsTimeout = setTimeout(function() {
        $('.alert-tips-div').fadeOut();
    }, options.time);
}


var DS_CONST_ICON_DEFAULT = 'https://nzjddxpun.qnssl.com/default.png';
var DS_CONST_LOADING_DEFAULT = 'http://7xp89t.com1.z0.glb.clouddn.com/dashboard-loading-default.gif';
var DS_CONST_THEME_DARK = '/assets/v2/images/theme-dark.png';
var DS_CONST_THEME_LIGHT = '/assets/v2/images/theme-light.png';
var DS_CONST_NOT_SET = '/assets/v2/images/not-set.jpg';
var DS_CONST_CUSTOM_DOWNLOAD_ANDROID =  '/assets/v2/images/custom-download-android.png';
var DS_CONST_CUSTOM_DOWNLOAD_IOS =  '/assets/v2/images/custom-download-ios.png';


module.exports = {
    initNav: initNav,
    renderUserInfo: renderUserInfo,
    renderApps: renderApps,
    getApps: getApps,
    selectApp: selectApp,

    uploadImage: uploadImage,
    showModal: showModal,
    hideModal: hideModal,

    showLoginModal: showLoginModal,
    hideLoginModal: hideLoginModal,

    showConfirmModal: showConfirmModal,
    showDeleteModal: showDeleteModal,

    initLoginModal: initLoginModal,

    alert: alertTips,

    DS_CONST_ICON_DEFAULT: DS_CONST_ICON_DEFAULT,
    DS_CONST_LOADING_DEFAULT: DS_CONST_LOADING_DEFAULT,
    DS_CONST_THEME_DARK: DS_CONST_THEME_DARK,
    DS_CONST_THEME_LIGHT: DS_CONST_THEME_LIGHT,
    DS_CONST_NOT_SET: DS_CONST_NOT_SET,
    DS_CONST_CUSTOM_DOWNLOAD_ANDROID: DS_CONST_CUSTOM_DOWNLOAD_ANDROID,
    DS_CONST_CUSTOM_DOWNLOAD_IOS: DS_CONST_CUSTOM_DOWNLOAD_IOS,
};

},{"../lib/util.js":3,"./api.js":1}],3:[function(require,module,exports){
function addDate(dadd, gran) {
    var a = new Date();
    a = a.valueOf();
    var b;
    var month;
    switch (gran) {
        case "d":
            a = a + dadd * 24 * 60 * 60 * 1000;
            a = new Date(a);
            month = a.getMonth() + 1;
            b = a.getDate() + "/" + month;
            break;
        case "w":
            a = a + dadd * 24 * 60 * 60 * 1000 * 7;
            a = new Date(a);
            month = a.getMonth() + 1;
            b = a.getDate() + "/" + month;
            break;
        case "m":
            a = new Date(a);
            month = (a.getMonth() + dadd + 12) % 12 + 1;
            b = month + "月";
            break;
    }
    return b;
}

function randomInt(max) {
    return 0;
    /*
    max = max || 1000;
    return Math.floor(Math.random() * max);
    */
}

function escapeHtml(text) {
    return text
        .replace(/&/g, '&amp;')
        .replace(/ /g, '&nbsp;')
        .replace(/\"/g, '&quot;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
}

function isFullUrl(url) {
    if (url && /^http/.test(url)) {
        return true;
    } else {
        return false;
    }
}

module.exports = {
    addDate: addDate,
    randomInt: randomInt,
    escapeHtml: escapeHtml,
    isFullUrl: isFullUrl,
};

},{}],4:[function(require,module,exports){
var DSCOMMON = require('../lib/common.js');


var sendEmail = function() {
  var stopPoint = disabelLink(60);
  $.post('/this-is-a-clandestine-resource/resent-email', {
          username  :   emailpath
  }, function(result) {
      if(!result.success) {
          window.clearInterval(stopPoint);
          $('.resend-div').children().remove()
          addResendLink();
      }
      DSCOMMON.showModal('邮件', result.message);
  }).fail(function(err) {
      window.clearInterval(stopPoint);
      $('.resend-div').children().remove()
      addResendLink();
      DSCOMMON.showModal('邮件', "邮件发送失败");
  });
}



var addResendLink = function() {
  $('.resend-div').append($("<a href='#' class='resend-email'>").html("重新发送"))
  $('.resend-email').on('click', sendEmail);
}

addResendLink();

var disabelLink = function(second) {
  $('.resend-div').children().remove()
  $('.resend-div').html($("<a>").html(second.toString() + "秒后重发"))
  var stopPoint = setInterval(function(){
    second -= 1;
    $('.resend-div').html($("<a>").html(second.toString() + "秒后重发"))
    if (second <= 0) {
      window.clearInterval(stopPoint);
      $('.resend-div').children().remove()
      addResendLink()
    }
  }, 1000);
  return stopPoint;
}
},{"../lib/common.js":2}]},{},[4])