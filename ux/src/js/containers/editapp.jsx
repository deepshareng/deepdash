var DSAPI = require('../lib/api.js');
var DSAPP = require('../lib/app.js');
var DSCOMMON = require('../lib/common.js');

var React = require('react');
var connect = require('react-redux').connect;

var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    // mapStateToProps must contain the full data
    var rst = {confInfo: state.confInfo};
    return Object.assign(rst, {
        themes: {
            light: '0',
            dark: '1', 
        },
    });
};

var mapDispatchToProps = function(dispatch) {
    return {
        updateConf: function(conf) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: conf, 
            });
        },
        onChangeAppName: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    appName: e.target.value
                }, 
            });
        },
        onChangeYybUrl: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    yyburl: e.target.value
                }, 
            });
        },
        onThemeChange: function(e) {
            // tricky: radio: onClick is used as onChange is react
            // and the clicked option is the activated option!
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    theme: e.target.value,
                }, 
            });
        },
        onChangeDownloadTitle: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    downloadTitle: e.target.value,
                }, 
            });
        },
        onChangeDownloadMsg: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    downloadMsg: e.target.value,
                }, 
            });
        },
        onChangeBundleId: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    iosBundler: e.target.value, 
                },
            });                 
        },
        onChangeIosLink: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    iosDownloadUrl: e.target.value, 
                },
            });                 
        },
        onChangeTeamId: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    iosTeamID: e.target.value, 
                },
            });                 
        },
        onChangeIos9OldYybEnable: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    iosYYBEnableBelow9: e.target.value, 
                },
            });                 
        },
        onChangeIos9NewYybEnable: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    iosYYBEnableAbove9: e.target.value, 
                },
            });                 
        },
        onChangeForceDownload: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    forceDownload: e.target.value, 
                },
            });                 
        },
        onChangePkgName: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    androidPkgname: e.target.value, 
                    androidHost: e.target.value,
                },
            });                 
        },
        onChangeAndroidYybEnable: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    androidYYBEnable: e.target.value, 
                },
            });                 
        },
        onChangeIsAndroidDownDirectly: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    androidIsDownloadDirectly: e.target.value, 
                },
            });                 
        },
        onChangeAndroidLink: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    androidDownloadUrl: e.target.value, 
                },
            });                 
        },
        onChangeSha256val: function(e) {
            dispatch({
                type: 'UPDATE_APP_CONF',
                conf: {
                    androidSHA256: e.target.value, 
                },
            });                 
        },
        
        initializeAppConf: function() {
            DSAPI.getAppInfo(sessionAppId, function(data) {
                // TODO: check 
                dispatch({
                    type: 'UPDATE_APP_CONF',
                    conf: {
                        appid: sessionAppId,
                        appName: data.appName,
                        yyburl: data.yyburl,
                        theme: data.theme || '0',
                        downloadTitle: data.downloadTitle,
                        downloadMsg: data.downloadMsg,
                        iconUrl: data.iconUrl || DSCOMMON.DS_CONST_ICON_DEFAULT,

                        userConfBgWeChatAndroidTipUrl: data.userConfBgWeChatAndroidTipUrl,
                        userConfBgWeChatIosTipUrl: data.userConfBgWeChatIosTipUrl,

                        iosBundler: data.iosBundler,
                        iosScheme: 'ds' + sessionAppId,
                        iosDownloadUrl: data.iosDownloadUrl,
                        iosTeamID: data.iosTeamID,

                        iosYYBEnableAbove9: data.iosYYBEnableAbove9 || 'false',
                        iosYYBEnableBelow9: data.iosYYBEnableBelow9 || 'false',
                        forceDownload: data.forceDownload || 'false',

                        androidPkgname: data.androidPkgname,
                        androidHost: data.androidPkgname,
                        androidYYBEnable: data.androidYYBEnable || 'false',

                        androidIsDownloadDirectly: data.androidIsDownloadDirectly, 
                        androidDownloadUrl: data.androidDownloadUrl,
                        androidSHA256: data.androidSHA256,
                        androidScheme: 'ds' + sessionAppId,
                    }, 
                });
                $('.yyb-link').trigger('change');  
                $('#isandroiddowndirectly').trigger('change');
            }, function() {
                DSCOMMON.alert('获取数据失败，请刷新页面');
            });


            $('#app-icon-react').on('change', function() {
                DSCOMMON.uploadImage(this, '#app-icon-show-area', '#app-icon-url', 'icon', function(data) {
                    dispatch({
                        type: 'UPDATE_APP_CONF',
                        conf: {
                            iconUrl: data.value,
                        }, 
                    });
                });
            });
            $('#custom-download-img-android-react').on('change', function() {
                DSCOMMON.uploadImage(this, '#custom-download-img-android-show-area', '#custom-download-img-android-url', 'conf', function(data) {
                    dispatch({
                        type: 'UPDATE_APP_CONF',
                        conf: {
                            userConfBgWeChatAndroidTipUrl: data.value,
                        }, 
                    });
                });
            });
            $('#custom-download-img-ios-react').on('change', function() {
                DSCOMMON.uploadImage(this, '#custom-download-img-ios-show-area', '#custom-download-img-ios-url', 'conf', function(data) {
                    dispatch({
                        type: 'UPDATE_APP_CONF',
                        conf: {
                            userConfBgWeChatIosTipUrl: data.value,
                        }, 
                    });
                });
            });
            $('#reset-download-ios-react').on('click', function() {
                dispatch({
                    type: 'UPDATE_APP_CONF',
                    conf: {
                        userConfBgWeChatIosTipUrl: '',
                    }, 
                });
            });
            $('#reset-download-android-react').on('click', function() {
                dispatch({
                    type: 'UPDATE_APP_CONF',
                    conf: {
                        userConfBgWeChatAndroidTipUrl: '',
                    }, 
                });
            });

            DSAPP.initAddApp();
        },
    };
};

var AppEditTpl = React.createClass({
    componentWillMount: function() {
    },
    componentDidMount: function() {
        this.props.initializeAppConf();
        DSGLOBAL.initComponent = this.props.initializeAppConf;
    },
    onClickComplete: function(e) {
        if (!DSAPP.checkLinkInput()) {
            return;
        }

        DSAPI.updateAppInfo(sessionAppId, this.props.confInfo, function() {
            DSAPI.setSessionAppId(sessionAppId);
            location.href = '#/appinfo'; 
        }, function(xhr) {
            if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_PARAM_INVALID) {
                DSCOMMON.showModal('创建应用', '创建应用失败，您输入的参数太长或不合法');
            } else {
                DSCOMMON.showModal('创建应用', '创建应用失败，请确认输入和网络状况，再重试');
            }
        });
    },
    onClickDelete: function() {
        DSCOMMON.showDeleteModal(function(){
        }, function(){
        }, function() {
            var confirmName = $('#confirmName').val();
            var deleteName = $('#appname').val(); 

            if (confirmName !== '' && confirmName === deleteName) {
                DSAPI.deleteApp(sessionAppId, function(){
                    $('#ds-delete-modal #info-msg').html('删除APP成功，将在2秒后跳转到首页...');
                    setTimeout(function(){
                        location.href = '/';
                    }, 2000);
                }, function() {
                    $('#ds-delete-modal #info-msg').html('删除APP失败');
                });
            } else {
                $('#ds-delete-modal #info-msg').html('输入的APP名称有误');
            }
        });
    },
    render: function() {
        return (
            <div className="app-info">
                <div className="row">
                  <div className="col-lg-12">
                    <h1 className="page-header">应用信息</h1>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                    <div className="panel panel-default app-info-panel">
                      <div className="panel-heading">
                        <h3 className="panel-title">通用配置</h3>
                      </div>
                      <div className="panel-body">
                        <form role="form">
                          <div className="form-group">
                            <label htmlFor="appname" className="control-label">应用名称</label>
                            <input id="appname" type="text" className="form-control" value={this.props.confInfo.appName} onChange={this.props.onChangeAppName}/>
                          </div>
                          <div className="form-group">
                            <label htmlFor="appid" className="control-label">appid</label>
                            <input id="appid" value={this.props.confInfo.appid} type="text" disabled className="form-control"/>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                    <div className="panel panel-default app-info-panel">
                      <div className="panel-heading">
                        <h3 className="panel-title">iOS应用基础配置</h3>
                      </div>
                      <div className="panel-body">
                        <form role="form">
                          <div id="bundleid-area" className="form-group">
                            <label htmlFor="bundleid" className="control-label">Bundle Identifier<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <input id="bundleid" type="text" className="form-control" value={this.props.confInfo.iosBundler} onChange={this.props.onChangeBundleId}/>
                          </div>
                          <div id="teamid-area" className="form-group">
                            <label htmlFor="teamid" className="control-label">Apple Team ID<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <input id="teamid" type="text" className="form-control" value={this.props.confInfo.iosTeamID} onChange={this.props.onChangeTeamId}/>
                          </div>
                          <div className="form-group">
                            <label htmlFor="ioslink" className="control-label">AppStore 下载地址</label>
                            <input id="ioslink" type="text" className="form-control" value={this.props.confInfo.iosDownloadUrl} onChange={this.props.onChangeIosLink}/>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                    <div className="panel panel-default app-info-panel">
                      <div className="panel-heading">
                        <h3 className="panel-title">iOS应用高级配置</h3>
                      </div>
                      <div className="panel-body">
                        <form role="form">
                          <div id="ios9-old-yyb-area" className="form-group">
                            <label htmlFor="ios9-old-yyb-enable" className="control-label">iOS9以下版本是否使用<a target="_blank" href="http://wiki.open.qq.com/index.php?title=mobile/应用宝微下载">应用宝微链接</a>下载链接<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <select value={this.props.confInfo.iosYYBEnableBelow9} id="ios9-old-yyb-enable" data-yyb-index="1" data-yyb-label="ios9old" className="form-control" onChange={this.props.onChangeIos9OldYybEnable}>
                              <option value="true">是</option>
                              <option value="false">否</option>
                            </select>
                          </div>
                          <div id="ios9-new-yyb-area" className="form-group">
                            <label htmlFor="ios9-new-yyb-enable" className="control-label">iOS9以上版本是否使用应用宝微下载链接<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <select value={this.props.confInfo.iosYYBEnableAbove9} id="ios9-new-yyb-enable" data-yyb-index="1" data-yyb-label="ios9new" className="form-control" onChange={this.props.onChangeIos9NewYybEnable}>
                              <option value="true">是</option>
                              <option value="false">否</option>
                            </select>
                          </div>
                          <div className="form-group">
                            <label htmlFor="force-download" className="control-label">强制下载</label>
                            <select value={this.props.confInfo.forceDownload} id="force-download" className="form-control" onChange={this.props.onChangeForceDownload}>
                              <option value="true">是</option>
                              <option value="false">否</option>
                            </select>
                          </div>
                          <div id="yyb-link-area-1" className="form-group hide yyb-link-area">
                            <label htmlFor="yyblink-1" className="control-label">应用宝微链接</label>
                            <input id="yyblink-1" type="text" className="form-control yyb-link" onChange={this.props.onChangeYybUrl} value={this.props.confInfo.yyburl}/>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                    <div className="panel panel-default app-info-panel">
                      <div className="panel-heading">
                        <h3 className="panel-title">Android应用基础信息</h3>
                      </div>
                      <div className="panel-body">
                        <form role="form">
                          <div id="pkgname-area" className="form-group">
                            <label htmlFor="pkgname" className="control-label">Package Name<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <input id="pkgname" type="text" className="form-control" value={this.props.confInfo.androidPkgname} onChange={this.props.onChangePkgName}/>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                    <div className="panel panel-default app-info-panel">
                      <div className="panel-heading">
                        <h3 className="panel-title">Android应用高级配置</h3>
                      </div>
                      <div className="panel-body">
                        <form role="form">
                          <div id="isandroiddowndirectly-area" className="form-group">
                            <label htmlFor="isandroiddowndirectly" className="control-label">通过浏览器打开的下载方式偏好<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <select id="isandroiddowndirectly" value={this.props.confInfo.androidIsDownloadDirectly} className="form-control" onChange={this.props.onChangeIsAndroidDownDirectly}>
                              <option value="true">直接下载APK </option>
                              <option value="false">应用市场 </option>
                            </select>
                            <p className="text-info help-block">推荐选择应用市场 </p>
                          </div>
                          <div id="android-link-area" className="form-group">
                            <label htmlFor="androidlink" className="control-label"><span>下载地址</span><i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <input id="androidlink" type="text" className="form-control" value={this.props.confInfo.androidDownloadUrl} onChange={this.props.onChangeAndroidLink}/>
                          </div>
                          <div id="android-yyb-enable-area" className="form-group">
                            <label htmlFor="android-yyb-enable" className="control-label">通过微信、QQ打开时启用应用宝微下载链接<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <select id="android-yyb-enable" value={this.props.confInfo.androidYYBEnable} data-yyb-index="2" data-yyb-label="android" className="form-control" onChange={this.props.onChangeAndroidYybEnable}>
                              <option value="true">是</option>
                              <option value="false">否</option>
                            </select>
                            <p className="text-info help-block">强烈建议启用应用宝微下载链接</p>
                          </div>
                          <div id="yyb-link-area-2" className="form-group hide yyb-link-area">
                            <label htmlFor="yyblink-2" className="control-label">应用宝微链接</label>
                            <input id="yyblink-2" type="text" className="form-control yyb-link" value={this.props.confInfo.yybUrl} onChange={this.props.onChangeYybUrl}/>
                          </div>
                          <div id="sha256val-area" className="form-group">
                            <label htmlFor="sha256val" className="control-label">SHA256签名指纹(可选)<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <input id="sha256val" type="text" className="form-control" value={this.props.confInfo.androidSHA256} onChange={this.props.onChangeSha256val}/>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div id="yyb-link-area-global" className="row">
                  <div className="col-lg-8">
                    <div className="panel panel-default app-info-panel">
                      <div className="panel-heading">
                        <h3 className="panel-title">通用高级配置</h3>
                      </div>
                      <div className="panel-body">
                        <form role="form">
                          <div className="form-group yyb-link-area">
                            <label htmlFor="yyblink-global" className="control-label">应用宝微链接</label>
                            <input id="yyblink-global" type="text" className="form-control yyb-link" value={this.props.confInfo.yybUrl} onChange={this.props.onChangeYybUrl}/>
                          </div>
                          <div className="form-group">
                            <label htmlFor="theme-option" className="control-label">主题</label>
                            <div>
                              <label className="radio-inline col-sm-4">
                                <input type="radio" onChange={this.props.onThemeChange} name="theme-option" value={this.props.themes.light} checked={this.props.confInfo.theme === this.props.themes.light } className="theme-option"/>light
                              </label>
                              <label className="radio-inline col-sm-4">
                                <input type="radio" onChange={this.props.onThemeChange} name="theme-option" value={this.props.themes.dark} checked={this.props.confInfo.theme === this.props.themes.dark} className="theme-option"/>dark
                              </label>
                            </div>
                            <div className="row">
                              <div className="col-sm-4"><img src="/assets/v2/images/theme-light.png" style={{ maxWidth:200, maxHeight: 360, }} className="theme-img-show-area"/></div>
                              <div className="col-sm-4"><img src="/assets/v2/images/theme-dark.png" style={{ maxWidth:200, maxHeight: 360, }} className="theme-img-show-area"/></div>
                            </div>
                          </div>
                          <div id="download-title-area" className="form-group">
                            <label htmlFor="download-title" className="control-label">下载提示信息1<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <p className="text-info help-block">可以在DeepShare Param里定制，使每次点击都有不同的提示信息</p>
                            <input id="download-title" type="text" className="form-control" value={this.props.confInfo.downloadTitle} onChange={this.props.onChangeDownloadTitle}/>
                          </div>
                          <div id="download-msg-area" className="form-group">
                            <label htmlFor="download-msg" className="control-label">下载提示信息2<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <input id="download-msg" type="text" className="form-control" value={this.props.confInfo.downloadMsg} onChange={this.props.onChangeDownloadMsg}/>
                          </div>
                          <div id="app-icon-area" className="form-group">
                            <label htmlFor="app-icon" style={{display: 'block'}} className="control-label">App图标<i className="fa fa-fw fa-lg fa-question-circle ds-tip"></i></label>
                            <p className="text-info help-block">会显示在提示下载页面中</p>
                            <img id="app-icon-show-area" src={this.props.confInfo.iconUrl || DSCOMMON.DS_CONST_ICON_DEFAULT}/>
                            <a style={{position: 'relative'}} className="btn btn-primary btn-outline pull-right">上传图标
                              <input id="app-icon-react" type="file" style={{position: 'absolute', top:0, left: 0, width: '100%', height: '100%', opacity: 0}}/></a>
                            <div style={{clear: 'both'}}></div>
                          </div>
                          <div className="form-group">
                            <label style={{display: 'block'}} className="control-label">自定义跳转页面(可选)</label>
                            <p className="text-info help-block">640X960 用于微信，qq当需要引导用户从Safari打开时</p>
                            <div className="row">
                              <div className="col-sm-4 download-images-area text-center">
                                <div style={{position: 'relative'}}>
                                <img id="custom-download-img-ios-show-area" className="custom-download-img-show-area" src={this.props.confInfo.userConfBgWeChatIosTipUrl || DSCOMMON.DS_CONST_CUSTOM_DOWNLOAD_IOS}/>
                                <img src="/assets/v2/images/download-ds-logo.png" className="custom-download-ds-logo"/></div>
                                <div><a style={{position: 'relative'}} className="btn btn-primary btn-outline">上传iOS版
                                    <input id="custom-download-img-ios-react" type="file" style={{position: 'absolute', top:0, left: 0, width: '100%', height: '100%', opacity: 0}}/></a><a id="reset-download-ios-react" style={{position: 'relative'}} className="btn btn-primary btn-outline">重置</a></div>
                              </div>
                              <div className="col-sm-4 download-images-area text-center">
                                <div style={{position: 'relative'}}>
                                <img id="custom-download-img-android-show-area" className="custom-download-img-show-area" src={this.props.confInfo.userConfBgWeChatAndroidTipUrl || DSCOMMON.DS_CONST_CUSTOM_DOWNLOAD_ANDROID}/>
                                <img src="/assets/v2/images/download-ds-logo.png" className="custom-download-ds-logo"/></div>
                                <div><a style={{position: 'relative'}} className="btn btn-primary btn-outline">上传Android版
                                    <input id="custom-download-img-android-react" type="file" style={{position: 'absolute', top:0, left: 0, width: '100%', height: '100%', opacity: 0}}/></a><a id="reset-download-android-react" style={{position: 'relative'}} className="btn btn-primary btn-outline">重置</a></div>
                              </div>
                              <div style={{clear: 'both'}}></div>
                            </div>
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                      <div className="pull-right">
                          <a id="submit-app-info" className="btn btn-primary btn-buttom-margin" onClick={this.onClickComplete}><i className="fa fa-fw fa-floppy-o"></i>保存修改</a>
                      <div style={{clear: 'both'}}></div>
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-8">
                      <div className="pull-right">
                          <a id="btn-delete-app" className="btn btn-danger btn-buttom-margin-last" onClick={this.onClickDelete}><i className="fa fa-fw fa-trash"></i>删除APP</a>
                      </div>
                      <div style={{clear: 'both'}}></div>
                  </div>
                </div>
                <div className="text-center ds-footer">
                  <footer>DeepShare Application &copy; 2016</footer>
                </div>
                <div id="ds-msg-modal" tabIndex="-1" role="dialog" aria-labelledby="ds-msg-modal-label" aria-hidden="true" style={{display: 'none'}} className="modal fade">
                  <div className="modal-dialog">
                    <div className="modal-content">
                      <div className="modal-header">
                        <button type="button" data-dismiss="modal" aria-hidden="true" className="close">×</button>
                        <h4 id="ds-msg-modal-label" className="modal-title">状态提示 </h4>
                      </div>
                      <div className="modal-body"></div>
                      <div className="modal-footer"></div>
                    </div>
                  </div>
                </div>
                <div id="ds-delete-modal" tabIndex="-1" role="dialog" aria-labelledby="ds-msg-modal-label" aria-hidden="true" style={{display: 'none'}} className="modal fade">
                  <div className="modal-dialog">
                    <div className="modal-content">
                      <div className="modal-header">
                        <button type="button" data-dismiss="modal" aria-hidden="true" className="close">×</button>
                        <h4 id="ds-msg-modal-label" className="modal-title">删除APP</h4>
                      </div>
                      <div className="modal-body">
                        <p className="text-danger">您将删除APP，相关数据不能恢复</p>
                        <div className="form-group">
                          <label htmlFor="confirmName" className="control-label">请输入此APP的名称，以确认删除</label>
                          <input id="confirmName" type="text" className="form-control"/>
                          <p id="info-msg" className="text-danger"></p>
                        </div>
                      </div>
                      <div className="modal-footer">
                        <div><a className="btn btn-default ds-cancel">取消</a><a className="btn btn-default ds-ensure">确认删除</a></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>);
    },
});

module.exports = {
    AppEditTpl: connect(mapStateToProps, mapDispatchToProps)(AppEditTpl)
};
;
