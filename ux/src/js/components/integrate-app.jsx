var React = require('react');
var ReactDOM = require('react-dom');
var connect = require('react-redux').connect;

var hljs = require('highlight.js');

var Link = require('react-router').Link;
var IndexLink = require('react-router').IndexLink;

var CompleteCreateTab = require('./integrate-snippet.jsx').CompleteCreateTab;
var CompleteConfTab = require('./integrate-snippet.jsx').CompleteConfTab;
var HorizontalProgressBar = require('./integrate-snippet.jsx').HorizontalProgressBar;
var CopyCodeBtn = require('./integrate-snippet.jsx').CopyCodeBtn;

var DSAPI = require('../lib/api.js');
var DSCOMMON = require('../lib/common.js');

var STYLE_TAB_ACTIVE = require('../styles/common.js').STYLE_TAB_ACTIVE;


var renderIosUnilinkCodeStr = function() {
    return '' + 
        '// Availability : iOS (9.0 and later)\n' +
        '- (BOOL)application:(UIApplication *)application continueUserActivity:(NSUserActivity *)userActivity restorationHandler:(void (^)(NSArray *))restorationHandler {\n' +
        '    //此处添加以下代码\n' +
        '    BOOL handledByDeepShare = [DeepShare continueUserActivity:userActivity];\n' +
        '    return handledByDeepShare;\n' +
        '}';
};

var renderIosSchemeCodeStrAbove = function() {
    return '' +
        '// Availability : iOS (9.0 and later)\n' +
        '- (BOOL)application:(UIApplication *)application openURL:(NSURL *)url options:(NSDictionary *)options {\n' +
        '    //此处添加以下代码\n' +
        '    if([DeepShare handleURL:url]){\n' +
        '        return YES;\n' +
        '    }\n' +
        '    return NO;\n' +
        '}\n';
};

var renderIosSchemeCodeStrBelow = function() {
    return '' +
        '// Availability : iOS (4.2 to 8.4)\n' +
        '- (BOOL)application:(UIApplication *)application openURL:(NSURL *)url sourceApplication:(nullable NSString *)sourceApplication annotation:(id)annotation {\n' +
        '    //此处添加以下代码\n' +
        '    if([DeepShare handleURL:url]){\n' +
        '        return YES;\n' +
        '    }\n' +
        '    return NO;\n' +
        '}\n';
};

var renderIosInitCodeStr = function(appid) {
    return '' +
        '- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions{\n' +
        '    // 请直接复制此行代码，粘贴到项目的目标位置\n' +
        '    [DeepShare initWithAppID:@"' + sessionAppId  + '" withLaunchOptions:launchOptions withDelegate:self];\n' +
        '    return YES;\n' +
        '}\n';
};

var renderIosInappDataCodeStr = function() {
    return '' +
        '- (void)onInappDataReturned: (NSDictionary *) params withError: (NSError *) error {\n' +
        '    if (!error) {\n' +
        '        NSLog(@"finished init with params = %@", [params description]);\n' +
        '        // *注意*: 下面一行替换成， 调用应用自己的接口跳转到分享时页面\n' +
        '        goToTarget(params);\n' +
        '    } else {\n' +
        '        NSLog(@"init error id: %ld %@",error.code, errorString);\n' +
        '    }\n' +
        '}\n';
};

// TODO: the jar file version
var renderAndroidDependencyCodeStr = function() {
    return '' +
        'dependencies {    \n' +
        '    compile files(\'libs/deepshare-v2.1.1.jar\')\n' +
        '}\n';
};

var renderAndroidManifestCodeStr = function() {
    return '' +
        '<manifest……>\n' +
        '    <uses-permission android:name="android.permission.INTERNET" />\n' +
        '    <application……>\n' +
        '        <activity……>\n' +
        '            <intent-filter>\n' +
        '                <data\n' +
        '                    android:host="此处填写DashBoard中显示的host"\n' +
        '                    android:scheme="此处填写DashBoard中显示的scheme" />\n' +
        '                <action android:name="android.intent.action.VIEW" />\n' +
        '                <category android:name="android.intent.category.DEFAULT" />\n' +
        '                <category android:name="android.intent.category.BROWSABLE" />\n' +
        '            </intent-filter>\n' +
        '        </activity>\n' +
        '    </application>\n' +
        '</manifest>\n';
};

var renderAndroidSetupCodeStr = function() {
    return  '' +
        '@Override\n' +
        'public void onNewIntent(Intent intent) {\n' +
        '    this.setIntent(intent);\n' +
        '}\n' +
        '\n' +
        '@Override\n' +
        'public void onStop() {\n' +
        '    super.onStop();\n' +
        '    DeepShare.onStop();//停止DeepShare\n' +
        '}\n';
};

var renderAndroidInappDataCodeStr = function(appid) {
    return '' +
        'public class MainActivity extends Activity implements DSInappDataListener {\n' +
        '    public void onStart() {\n' +
        '        super.onStart();\n' +
        '        // 请直接复制下一行代码，粘贴到相应位置\n' +
        '        DeepShare.init(this, "' + appid + '", this);\n' +
        '    }\n' +
        '\n' +
        '    @Override\n' +
        '        /** 代理方法onInappDataReturned处理获取的启动参数\n' +
        '         * @param params 所获取到的启动参数\n' +
        '         */\n' +
        '        public void onInappDataReturned(JSONObject params) {\n' +
        '            try {\n' +
        '                if (params == null) return;\n' +
        '                String cmdName = params.getString("name");\n' +
        '                // *注意*: 下面一行替换成， 调用应用自己的接口跳转到分享时页面\n' +
        '                goToTarget(params); \n' +
        '            } catch (JSONException e) {\n' +
        '                e.printStackTrace();\n' +
        '            }\n' +
        '        }\n' +
        '\n' +
        '    @Override\n' +
        '    public void onFailed(String s) {\n' +
        '        // handle failure\n' +
        '    }\n' +
        '}\n';
};

var IosTpl = React.createClass({
    render: function() {
        var scheme = hljs.highlight('bash', 'ds' + this.props.appid).value;
        var IosUnilinkCode = hljs.highlight('objectivec', renderIosUnilinkCodeStr()).value;
        var IosSchemeCodeAbove = hljs.highlight('objectivec', renderIosSchemeCodeStrAbove()).value;
        var IosSchemeCodeBelow = hljs.highlight('objectivec', renderIosSchemeCodeStrBelow()).value;
        var IosInitCode = hljs.highlight('objectivec', renderIosInitCodeStr(this.props.appid)).value;
        var IosInappDataCode = hljs.highlight('objectivec', renderIosInappDataCodeStr()).value;

        return ( 
            <div style={{marginTop: '2rem'}}>
            <ol>
                <li>
                    <h5>下载新版SDK</h5> 
                    <p>下载<a target="_blank" href="https://deepshare.blob.core.chinacloudapi.cn/sdk/deepshare_ios_v2.1.9.zip">iOS SDK</a></p>
                    <p>将lib文件夹的所有文件拖到工程目录</p>
                </li>
                <li>
                    <h5>配置*-info.plist文件</h5> 
                    <p>添加用户打开APP的URL Scheme</p>
                    <ol type="a">
                        <li>添加一个叫URL types的键值</li>
                        <li>点击左边剪头打开列表，可以看到Item 0，一个字典实体</li>
                        <li>点击Item 0新增一行，从下拉列表中选择URL Schemes，敲击键盘回车键完成插入</li>
                        <li>更改所插入URL Schemes的值为DashBoard中生成的Scheme(如下所示)</li>
                    </ol> 
                    <pre>
                        <code className="bash hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: scheme}}></div>
                        </code>
                    </pre>
                    <em className="small text-info">注:请直接复制此行中的scheme，并粘贴到相应位置</em>
                    <p>
                        完成后，如下图所示
                        <img src="/assets/v2/images/integrate-ios-scheme.png"/>
                    </p>
                </li>
                <li>
                    <h5>配置universal link</h5> 
                    <p>配置开发者信息</p>
                    <ol type="a">
                        <li>
                            登录developers.apple.com，点击按钮「Certificate, Identifiers & Profiles」,再点击「Identifiers」

                            <img src="/assets/v2/images/integrate-ios-developer-1.png"/>
                        </li>
                        <li>
                            确保开启「Associated Domains」，这个按钮在页面下方，如图所示
                            <img src="/assets/v2/images/integrate-ios-developer-2.png"/>
                        </li>
                    </ol>
                    <p>配置Xcode</p>
                    <ol type="a">
                        <li>
                            在Xcode配置中开启「Associated Domains」
                            <ol type="i">
                                <li> 选择相应的target </li>
                                <li> 点击「Capabilities tab」 </li>
                                <li> 开启Associated Domains </li>
                                <li> 点击「+」按钮，添加一个Associated Domain，其内容为applinks:fds.so </li>
                            </ol>
                            完成后，如图所示
                            <img src="/assets/v2/images/integrate-ios-developer-3.png"/>
                        </li>
                        <li>
                            在Xcode中添加依赖库
                            <ol type="i">
                                <li> 选择相应的target </li>
                                <li> 选择Build Phrases </li>
                                <li> 在Link Binary With Libraries中添加SafariServices.framework </li>
                            </ol>
                        </li>
                    </ol>
                    <p>添加以下代码到AppDelegate</p>
                    <pre>
                        <code className="objectivec hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: IosUnilinkCode}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <h5>添加调用代码</h5> 
                    <p>添加初始化代码</p>
                    <pre>
                        <code className="objectivec hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: IosSchemeCodeAbove}}></div>
                        </code>
                    </pre>
                    <p>如果是iOS8和以下版本，请用以下代码</p>
                    <pre>
                        <code className="objectivec hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: IosSchemeCodeBelow}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <h5>获取场景还原的代码</h5> 
                    <p>在AppDelegate的didFinishLaunchingWithOptions方法中添加[DeepShare initWithAppID…]的调用，并通过Delegate接受返回的启动参数</p>
                    <pre>
                        <code className="objectivec hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: IosInitCode }}></div>
                        </code>
                    </pre>
                    <p>获取场景还原的参数</p>
                    <pre>
                        <code className="objectivec hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: IosInappDataCode }}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <p>完成iOS SDK集成,  如有问题请参考<a target="_blank" href="http://deepshare.io/doc/#ios-sdk">开发文档</a></p>
                </li>
            </ol> 
            </div>
        );
    },
});

var AndroidTpl = React.createClass({
    render: function() {
        var AndroidDependencyCode = hljs.highlight('nginx', renderAndroidDependencyCodeStr()).value;
        var AndroidManifestCode = hljs.highlight('xml', renderAndroidManifestCodeStr()).value;
        var AndroidSetupCode = hljs.highlight('java', renderAndroidSetupCodeStr()).value;
        var AndroidInappDataCode = hljs.highlight('java', renderAndroidInappDataCodeStr(this.props.appid)).value;

        return ( 
            <div style={{marginTop: '2rem'}}>
            <ol>
                <li>
                    <h5>下载新版SDK</h5>
                    <p>下载<a target="_blank" href="https://deepshare.blob.core.chinacloudapi.cn/sdk/deepshare_android_v2.1.1.zip">Android SDK</a></p>
                    <p>将deepshare-vx.x.x.jar 添加到工程的libs目录下，并重新build工程</p>
                    <p>
                    Eclipse ADT: 解压下载的SDK目录，并将libs目录下的deepshare-vx.x.x.jar拷贝到工程目录的libs文件夹下
                    </p>
                    <p>
                    Android Studio: 
                    <ul>
                        <li>选择project视图，将deepshare-vx.x.x.jar添加到app目录下的libs目录中（如不存在libs目录，请新建libs目录再添加）</li>
                        <li>
                            <p>在app目录的build.gradle文件中添加dependencies依赖项</p>
                            <pre>
                                <code className="nginx hljs">
                                <CopyCodeBtn />
                                    <div dangerouslySetInnerHTML={{__html: AndroidDependencyCode}}></div>
                                </code>
                            </pre>
                        </li>
                    </ul>
                    </p>
                </li>
                <li>
                    <h5>修改Manifext.xml</h5> 
                    <p>在启动Activity中增加DeepShare的intent-filter，这样您的App就可以通过浏览器被唤起</p>
                    <pre>
                        <code className="xml hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: AndroidManifestCode}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <h5>启动和停止的回调</h5> 
                    <p>在Activity的<code>onNewIntent</code>和<code>onStop</code>添加如下代码</p>
                    <pre>
                        <code className="java hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: AndroidSetupCode}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <h5>添加场景还原代码</h5> 
                    <p>在<strong className="text-info">启动Activity</strong>添加以下代码</p>
                    <pre>
                        <code className="java hljs">
                        <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: AndroidInappDataCode}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <p>完成Android SDK集成,  如有问题请参考<a href="http://deepshare.io/doc/#android-sdk">开发文档</a></p>
                </li>
            </ol> 
            </div>
        );
    },
});


var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    return {
        appInfo: state.appIntegrateInfo, 
    };
};

var mapDispatchToProps = function(dispatch) {
    return {
        switchAppId: function(appid) {
            dispatch({
                type: 'UPDATE_APP_ID',
                appid: appid,
            });
        },
    };
};
var AppTpl = React.createClass({
    getInitialState: function() {
        return {};
    },
    onClickNext: function() {
        location.href = '#/web';
    },
    onClickPrev: function() {
        location.href = '#/configure';
    },
    componentDidMount: function() {
        var self = this;
        DSAPI.getAppInfo(sessionAppId, function(data) {
            self.props.switchAppId(sessionAppId);
        }, function(xhr) {
            if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                DSCOMMON.showLoginModal()
            }
        });

        DSGLOBAL.initComponent = function() {
            DSAPI.getAppInfo(sessionAppId, function(data) {
                self.props.switchAppId(sessionAppId);
            }, function(xhr) {
                if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                    DSCOMMON.showLoginModal()
                }
            });
        };
    },
    componentDidUpdate: function() {
        $('body').animate({
            scrollTop: 0,
        }, 150);
    },
    render: function() {
        return (
         <div>
            <HorizontalProgressBar stepIndex={3}/>   

          <div className="row ds-integrate-steps-content">
            <div className="col-lg-8 col-lg-offset-2">
              <div className="panel panel-default">
                <div className="panel-heading">
                  <h3 className="panel-title"> 
                      <span className="btn btn-info btn-circle ds-integrate-status-icon">
                          <i className="fa fa-arrow-right"></i>
                      </span>
                      <span>3.APP中集成SDK</span>
                  </h3>
                </div>
                <div className="panel-body">
                    <ul role="nav" className="nav nav-tabs">
                        <li><IndexLink to="/app/ios" activeStyle={STYLE_TAB_ACTIVE}>iOS</IndexLink></li>
                        <li><Link to="/app/android" activeStyle={STYLE_TAB_ACTIVE}>Android</Link></li>
                    </ul>
                    {this.props.children && React.cloneElement(this.props.children, {appid: this.props.appInfo.appid})}
                </div>
                <div className="panel-footer">
                  <div className="pull-left"><a onClick={this.onClickPrev} className="btn btn-primary">上一步</a></div>
                  <div className="pull-right"><a onClick={this.onClickNext} className="btn btn-primary">下一步</a></div>
                  <div style={{clear: 'both'}}></div>
                </div>
              </div>
            </div>
          </div>
        </div>   
        );        
    },

});

// TODO: only pass needed props
//{this.props.children && React.cloneElement(this.props.children, this.props)}

module.exports = {
    IosTpl: IosTpl,
    AndroidTpl: AndroidTpl,
    AppTpl: connect(mapStateToProps, mapDispatchToProps)(AppTpl),
};
