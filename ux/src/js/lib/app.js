var DSUTIL = require('../lib/util.js');
var DSCOMMON = require('../lib/common.js');

var yybEnable = {
    android: false,
    ios9old: false,
    ios9new: false,
};
var yybLink = '';
var yybEnableCount = 0;

function checkLinkInput() {
    if (yybLink && !DSUTIL.isFullUrl(yybLink)) {
        DSCOMMON.showModal('字段错误', '请确保应用宝链接完整，包含协议头: http://或https://');
        return false;
    }
    /* remove for default value: '-'
    if ($('#androidlink').val() && !DSUTIL.isFullUrl($('#androidlink').val())) {
        DSCOMMON.showModal('字段错误', '请确保Android下载地址完整，包含协议头: http://或https://');
        return false;
    }*/
    if ($('#ioslink').val() && !DSUTIL.isFullUrl($('#ioslink').val())) {
        DSCOMMON.showModal('字段错误', '请确保iOS下载地址完整，包含协议头: http://或https://');
        return false;
    }

    return true;
}

function checkYyb() {
    if (yybEnableCount > 0 && yybLink === '') {
        return false;
    } else {
        return true;
    } 
}

function showInputError(selector, timeSpan) {
    $(selector).addClass('has-error');
    setTimeout(function() {
        $(selector).removeClass('has-error');
    }, timeSpan);
}

function showInputSuccess(selector, timeSpan) {
    $(selector).addClass('has-success');
    setTimeout(function() {
        $(selector).removeClass('has-success');
    }, timeSpan);
}

function showYybLink(enable, label, index) {
    if($('#yyb-link-area-global').length > 0) {
       if(enable === 'true') {
            if(!yybLink ) {
                window.location = '#yyb-link-area-global';
                showInputError('.yyb-link-area', 3000);
            } else {
                //showInputSuccess('.yyb-link-area', 3000);
            }
        }
    } else {
        if(enable === 'true') {
            if (yybEnable[label] === false) {
                yybEnable[label] = true;
                yybEnableCount = yybEnableCount + 1;
                if (yybEnableCount === 1) {
                    $('.yyb-link-area').addClass('hide');
                    $('#yyb-link-area-' + index).removeClass('hide');
                }
            }
        } else if (enable === 'false') {
            if (yybEnable[label] === true) {
                yybEnable[label] = false;
                yybEnableCount = yybEnableCount - 1;
                if (yybEnableCount === 0) {
                    $('.yyb-link-area').addClass('hide');
                }
            }
        } 
    }
}

function initAddApp() {
    yybEnable = {
        android: false,
        ios9old: false,
        ios9new: false,
    };
    yybLink = '';
    yybEnableCount = 0;

    $('#app-icon').on('change', function() {
        DSCOMMON.uploadImage(this, '#app-icon-show-area', '#app-icon-url', 'icon');
    });
    $('#custom-download-img-android').on('change', function() {
        DSCOMMON.uploadImage(this, '#custom-download-img-android-show-area', '#custom-download-img-android-url', 'conf');
    });
    $('#custom-download-img-ios').on('change', function() {
        DSCOMMON.uploadImage(this, '#custom-download-img-ios-show-area', '#custom-download-img-ios-url', 'conf');
    });
    $('#reset-download-ios').on('click', function() {
        $('#custom-download-img-ios-url').val(''); 
        $('#custom-download-img-ios-show-area').attr('src', DSCOMMON.DS_CONST_CUSTOM_DOWNLOAD_IOS);
    });
    $('#reset-download-android').on('click', function() {
        $('#custom-download-img-android-url').val(''); 
        $('#custom-download-img-android-show-area').attr('src', DSCOMMON.DS_CONST_CUSTOM_DOWNLOAD_ANDROID);
    });

    $('#isandroiddowndirectly').on('change', function() {
        if($(this).val() === 'true') {
            // TODO:
            $('#android-link-area label span').html('下载地址');
            $('#android-link-area label .ds-tip').addClass('hide');
            //$('#android-link-area').removeClass('hide');
        } else if ($(this).val() === 'false') {
            // TODO:
            $('#android-link-area label span').html('Chrome应用下载地址');
            $('#android-link-area label .ds-tip').removeClass('hide');
            //$('#android-link-area').addClass('hide');
        } 
    });
    $('#ios9-old-yyb-enable, #ios9-new-yyb-enable, #android-yyb-enable').on('change', function() {
        var enable = $(this).val();
        var index = $(this).data('yyb-index');
        var label = $(this).data('yyb-label');
        showYybLink(enable, label, index);
    });
    $('.yyb-link').on('change', function() {
        $('.yyb-link-area').removeClass('has-error');
        yybLink = $(this).val().trim();
        $('.yyb-link').val(yybLink);
    });

    $("#bundleid-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '<img src="/assets/v2/images/tip-bundle-identifer.png" width="450">',
    });
    $("#download-title-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '<img src="/assets/v2/images/download-title.jpg" width="250" height="450">',
    });
    $("#download-msg-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '<img src="/assets/v2/images/download-msg.jpg" width="250" height="450">',
    });
    $("#download-icon-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '<img src="/assets/v2/images/download-icon.jpg" width="250" height="450">',
    });
    $("#yybselect-area .ds-tip").popover({
        toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '启用应用宝微链接，从微信中跳转可以省略经过safari的一步',
    });
    $("#teamid-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '请参考Apple官网 develop.apple.com',
        content: '<img src="/assets/v2/images/teamid.jpg" width="450" height="350">',
    });
    $("#ios9-old-yyb-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '使用应用宝微下载链接时用户会通过应用宝页面跳转到商店，再从商店打开App。<br>不使用应用宝微下载链接时用户会被引导点击右上角菜单, 并选择打开Safari, 之后通过Safari自动打开App',
    });
    $("#ios9-new-yyb-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: 'iOS 9以上OS版本通过Safari跳转时，新应用安装时可以做到100%匹配。<br>通过应用宝微下载链接跳转时，新应用安装时匹配精确度有所下降',
    });
    $("#android-yyb-enable-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '启用应用宝微下载链接，可以避免用户在微信，QQ中手动选择通过浏览器打开。<br>可以提高下载转化率',
    });
    $("#pkgname-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '<img src="/assets/v2/images/tip-package-name.png" width="450">',
    });
    $("#isandroiddowndirectly-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '启用“直接下载APK”则应用未安装时通过浏览器直接下载APK，某些浏览器可能不支持直接下载<br/>启用“跳转商店”则应用未安装时打开应用商店',
    });
    $("#sha256val-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '如何获取 Sha256_Cert_Fingerprints',
        content: '填写后可以自动激活Android M的applink功能<br>命令行执行此命令 : echo | keytool -list -v -keystore /path/to/app/release-key.keystore 2> /dev/null | grep SHA256',
    });
    $("#app-icon-area .ds-tip").popover({
        html: true, toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '<img src="/assets/v2/images/download-icon.jpg" width="250" height="450">',
    });
    $("#android-link-area .ds-tip").popover({
        toggle: 'popover', trigger: 'hover', placement: 'right', title: '',
        content: '当前浏览器是Chrome并且当应用未安装时，需要制定一个回调地址，通常是App的下载地址',
    });
}

module.exports = {
    initAddApp: initAddApp,
    checkLinkInput: checkLinkInput,
    checkYyb: checkYyb,
    yybEnable: yybEnable,
    yybLink: function() { return yybLink; },
    showInputError: showInputError,
    showInputSuccess: showInputSuccess,
};
