var DSCOMMON = require('../lib/common.js');
var DSAPI = require('../lib/api.js');

var React = require('react');
var connect = require('react-redux').connect;
var Panel = require('../components/panel.jsx')
var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    return {
        appInfo: state.appInfo, 
    };
};

var mapDispatchToProps = function(dispatch) {
    return {
        initializeAppInfo: function() {
            DSAPI.getAppInfo(sessionAppId, function(data) {
                dispatch({
                    type: 'INIT_APP_INFO',
                    appInfo: data,
                });
            }, function() {
                DSCOMMON.alert('获取数据失败，请刷新页面');
            });
        },
    };
};

var themeDict = {
    '0': {'title': '浅色', 'demo': DSCOMMON.DS_CONST_THEME_LIGHT,},
    '1': {'title': '深色', 'demo': DSCOMMON.DS_CONST_THEME_DARK,},
};
var AppInfoTpl = React.createClass({
    componentWillMount: function() {
        this.props.initializeAppInfo();
        DSGLOBAL.initComponent = this.props.initializeAppInfo;
    },
    componentDidMount: function() {
    },
    onClickEdit: function(e) {
        location.href = '#editapp';
    },
    basicInfo: function() {
      return [{title: "应用名称", content: this.props.appInfo.appName, id: "appname"},
              {title: "appid", content: sessionAppId, id: "appid"}];
    },
    iosInfo: function() {
      return [{title: "Bundle Identifier", content: this.props.appInfo.iosBundler || '未设置', id: "bundleid"},
              {title: "Scheme",            content: this.props.appInfo.iosScheme || 'ds' + sessionAppId, id: "iosscheme"},
              {title: "AppStore 下载地址",  content: this.props.appInfo.iosDownloadUrl || '未设置', id: "ioslink"},
              {title: "Apple账号 Team ID",  content: this.props.appInfo.iosTeamID || '未设置', id: "teamid"}];
    },
             
    iosConfig: function() {
      return [
          {title: "iOS9以下启用应用宝链接", content: (this.props.appInfo.iosYYBEnableBelow9 === 'true') && '是' || '否', id:"ios9-old-yyb-enable"},
          {title: "iOS9以上启用应用宝链接", content: (this.props.appInfo.iosYYBEnableAbove9 === 'true') && '是' || '否', id:"ios9-new-yyb-enable"},
          {title: "强制跳转应用宝", content: (this.props.appInfo.forceDownload === 'true') && '是' || '否', id:"force-download"},
      ];
    
    },
    andriodInfo: function() {
      return [{title: "Package Name", id: "pkgname", content: this.props.appInfo.androidPkgname || '未设置'},
              {title: "Scheme", id: "adrdscheme", content: this.props.appInfo.androidScheme || 'ds' + sessionAppId},
              {title: "Host", id: "androidhost", content: this.props.appInfo.androidHost || '未设置(请设置package name)'}]
    },
    andriodConfig: function() {
      return [{title: "下载方式", id: "androidisdownloaddirectly", content: (this.props.appInfo.androidIsDownloadDirectly === 'true') && '直接下载' || '商店下载'},
              {title: "下载地址", id: "androidlink", content: this.props.appInfo.androidDownloadUrl || '未设置'},
              {title: "SHA256", id: "sha256val", content: this.props.appInfo.androidSHA256 || '未设置'},
              {title: "Android启用应用宝链接", id:"android-yyb-enable", content: (this.props.appInfo.androidYYBEnable === 'true') && '是' || '否'}]
    },
    advancedConfig: function() {
      var themeTitle = themeDict[this.props.appInfo.theme || '0'].title;
      var themeDemo = themeDict[this.props.appInfo.theme || '0'].demo;
      var result = [];
      result.push({title: "应用宝微链接", id:"yyblink", content:this.props.appInfo.yyburl || '未设置'});
      result.push({
        title: "主题",
        content: [<span id="theme-option">{themeTitle}</span>,<br/>,
                  <img id="theme-img" src={themeDemo} className="theme-img-show-area"/>]
        });
      result.push({title: "下载提示信息1", id: "download-title" , content: this.props.appInfo.downloadTitle || '未设置'});
      result.push({title: "下载提示信息2", id: "download-msg", content: this.props.appInfo.downloadMsg || '未设置'});
      result.push({
        title: "App图标", 
        content: <img id="app-icon-show-area" src={this.props.appInfo.iconUrl || DSCOMMON.DS_CONST_ICON_DEFAULT}/>
      });
      result.push({
        title: "自定义跳转页面", 
        content:
          <div className="row">
            <div className="col-sm-5 download-images-area text-center">
              <div style={{'position': 'relative'}}>
                  <img src={this.props.appInfo.userConfBgWeChatIosTipUrl || DSCOMMON.DS_CONST_CUSTOM_DOWNLOAD_IOS} id="custom-download-img-ios-show-area" className="custom-download-img-show-area"/>
                  <img src="/assets/v2/images/download-ds-logo.png" className="custom-download-ds-logo"/></div>
                  <p>iOS</p>
            </div>
            <div className="col-sm-5 download-images-area text-center">
              <div style={{'position': 'relative'}}>
                  <img src={this.props.appInfo.userConfBgWeChatAndroidTipUrl || DSCOMMON.DS_CONST_CUSTOM_DOWNLOAD_ANDROID} id="custom-download-img-android-show-area" className="custom-download-img-show-area"/>
                  <img src="/assets/v2/images/download-ds-logo.png" className="custom-download-ds-logo"/></div>
                  <p>Android</p>
            </div>
            <div style={{clear: 'both'}}></div>
          </div>
      });
      return result;
    },
    render: function() {
        return (
            <div className="app-info">
                <div className="row">
                  <div className="col-lg-12">
                    <h1 className="page-header">应用信息</h1>
                  </div>
                </div>
                <Panel.row size="col-lg-8" theme="app-info-panel" title="基础信息" content={this.basicInfo()}/>
                <Panel.row size="col-lg-8" theme="app-info-panel" title="iOS应用信息" content={this.iosInfo()}/>
                <Panel.row size="col-lg-8" theme="app-info-panel" title="iOS应用高级配置" content={this.iosConfig()}/>
                <Panel.row size="col-lg-8" theme="app-info-panel" title="Android应用信息" content={this.andriodInfo()}/>
                <Panel.row size="col-lg-8" theme="app-info-panel" title="Android应用高级配置" content={this.andriodConfig()}/>
                <Panel.row size="col-lg-8" theme="app-info-panel" title="高级配置" content={this.advancedConfig()}/>
                <div className="row">
                  <div className="col-lg-8">
                      <a id="btn-edit-app-info" className="btn btn-primary pull-right btn-buttom-margin-last" onClick={this.onClickEdit}>
                          <i className="fa fa-fw fa-pencil-square-o"></i>
                          更新应用信息
                      </a>
                  </div>
                </div>
              </div>);        
    },
});

module.exports = {
    AppInfoTpl: connect(mapStateToProps, mapDispatchToProps)(AppInfoTpl)
};
