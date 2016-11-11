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
