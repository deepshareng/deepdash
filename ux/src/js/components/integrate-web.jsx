var React = require('react');
var ReactDOM = require('react-dom');
var connect = require('react-redux').connect;

var hljs = require('highlight.js');

var Link = require('react-router').Link;
var IndexLink = require('react-router').IndexLink;

var CompleteCreateTab = require('./integrate-snippet.jsx').CompleteCreateTab;
var CompleteConfTab = require('./integrate-snippet.jsx').CompleteConfTab;
var CompleteAppTab = require('./integrate-snippet.jsx').CompleteAppTab;
var HorizontalProgressBar = require('./integrate-snippet.jsx').HorizontalProgressBar;
var CopyCodeBtn = require('./integrate-snippet.jsx').CopyCodeBtn;

var DSAPI = require('../lib/api.js');
var STYLE_TAB_ACTIVE = require('../styles/common.js').STYLE_TAB_ACTIVE;

var renderStartRedirectCodeStr = function(appid) {
    return '' +
    '<script type="text/javascript">\n' +
    '// 组装参数\n' +
    'var params = {\n' +
    '    inapp_data: {\n' +
    '        name: \'mkdir\',\n' +
    '    }\n' +
    '};\n' +
    '\n' +
    '// 初始化DeepShare\n' +
    'var deepshare = new DeepShare(\'' + appid + '\');\n' +
    'deepshare.BindParams(params);\n' +
    '\n' +
    'document.getElementById(\'openOrDownloadApp\').addEventListener(\'click\', function() {\n' +
    '    // 触发跳转\n' +
    '    deepshare.Start();\n' +
    '});\n' +
    '</script>\n';
};

var renderWebApiCodeStr = function(appid) {
    return '' +
        'var params = {inapp_data: {name: \'mkdir\'}};\n' +
        '$.ajax({\n' +
        '    url: \'https://fds.so/v2/url/' + appid + '\',\n' +
        '    type: \'POST\',\n' +
        '    data: JSON.stringify(params),\n' +
        '    xhrFields: {withCredentials: true,},\n' +
        '    success: function(result) {\n' +
        '        //deepshareUrl为所生成的DeepShare Link, 可以让用户跳转到此链接 打开/下载 APP \n' +
        '        deepshareUrl = result.url;\n' +
        '    },\n' +
        '    error: function() {\n' +
        '        //出错情况下可以继续用户正常逻辑，可设置一个默认url.\n' +
        '    },\n' +
        '});\n';
};

var renderWebApiExampleCodeStr = function(appid) {
    return '' +
        '<script>\n' +
        'var deepshareUrl = null;\n' +
        'var params = {inapp_data: {name: \'mkdir\'}};\n' +
        '$.ajax({\n' +
        '    url: \'https://fds.so/v2/url/' + appid + '\',\n' +
        '    type: \'POST\',\n' +
        '    data: JSON.stringify(params),\n' +
        '    xhrFields: {withCredentials: true,},\n' +
        '    success: function(result) {\n' +
        '        //deepshareUrl为所生成的DeepShare Link, 可以让用户跳转到此链接 打开/下载 APP \n' +
        '        deepshareUrl = result.url;\n' +
        '    },\n' +
        '    error: function() {\n' +
        '        //出错情况下可以继续用户正常逻辑，可设置一个默认url.\n' +
        '    },\n' +
        '});\n' +

        'window.onload = function() {\n' +
        '    // 请将 #openOrDownloadApp 替换成页面中用户点击的按钮\n' +
        '    $(\'#openOrDownloadApp\').on(\'click\', function(){\n' +
        '        if(deepshareUrl != null){\n' +
        '            location.href = deepshareUrl;\n' +
        '        }else{\n' +
        '            //出错情况下可以继续用户正常逻辑.例如可以设置直接前往APP商店\n' +
        '        } \n' +
        '    });\n' +
        '};\n' +
        '</script>\n';
};

var ApiTpl = React.createClass({
    render: function() {
        var codeSnippet = renderWebApiCodeStr(this.props.appid);
        var highlightedCode = hljs.highlight('javascript', codeSnippet).value;

        var codeSnippetExample = renderWebApiExampleCodeStr(this.props.appid);
        var highlightedCodeExample = hljs.highlight('html', codeSnippetExample).value;
        return (
            <ol>
                <li>
                    <h5>调用DeepShare接口，获取跳转链接</h5>
                    <pre>
                        <code className="javascript hljs">
                        <CopyCodeBtn />
                        <div dangerouslySetInnerHTML={{__html: highlightedCode}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <h5>下面集成跳转的完整例子</h5>
                    <pre>
                        <code className="html hljs">
                            <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: highlightedCodeExample}}></div>
                        </code>
                    </pre>
                    <em className="small text-info">注:请直接复制以上代码到所开发的网页中body元素的最后，并将#openOrDownloadApp替换成自己页面中的按钮ID</em>
                    <br/>
                    <em className="small text-info">注:以上代码依赖jQuery，请确保页面中引用了jQuery</em>
                </li>
            </ol>  
        );        
    },
});

var JsTpl = React.createClass({
    render: function() {
        var highlightedCode = hljs.highlight('html', renderStartRedirectCodeStr(this.props.appid)).value;
        var highlightedCodeImport = hljs.highlight('html', '<script src="http://qn.fds.so/deepshare_v2.7.min.js"></script>').value;
        return (
            <ol>
                <li>
                    <h5>在页面中引入DeepShare的js</h5>
                    <pre>
                        <code className="html hljs">
                            <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: highlightedCodeImport}}></div>
                        </code>
                    </pre>
                </li>
                <li>
                    <h5>调用进行跳转</h5>
                    <pre>
                        <code className="html hljs">
                            <CopyCodeBtn />
                            <div dangerouslySetInnerHTML={{__html: highlightedCode}}></div>
                        </code>
                    </pre>
                    <em className="small text-info">注:请直接复制以上代码到所开发的网页中body元素的最后，并将#openOrDownloadApp替换成自己页面中的按钮ID</em>
                    <br/>
                    <em className="small text-info">注:2.6以前版本依赖jQuery，请确保页面中引用了jQuery。2.7以后版本不再需要依赖jQuery</em>
                </li>
            </ol>  
        );        
    },
});

var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    return {
        webInfo: state.webInfo, 
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

var WebTpl = React.createClass({
    onClickNext: function() {
        location.href = '#/complete';
    },
    onClickPrev: function() {
        location.href = '#/app';
    },
    componentDidMount: function() {
        var self = this;
        DSAPI.getAppInfo(sessionAppId, function(data) {
            self.props.switchAppId(sessionAppId);
        });
        
        DSGLOBAL.initComponent = function() {
            DSAPI.getAppInfo(sessionAppId, function(data) {
                self.props.switchAppId(sessionAppId);
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
            <HorizontalProgressBar stepIndex={4}/>

          <div className="row ds-integrate-steps-content">
            <div className="col-lg-8 col-lg-offset-2">
              <div className="panel panel-default">
                <div className="panel-heading">
                  <h3 className="panel-title"> <span className="btn btn-info btn-circle ds-integrate-status-icon"><i className="fa fa-arrow-right"></i></span><span>4.网页中集成deepshare</span></h3>
                </div>
                <div className="panel-body">
                    <ul role="nav" className="nav nav-tabs">
                        <li><IndexLink to="/web/api" activeStyle={STYLE_TAB_ACTIVE}>Web API</IndexLink></li>
                        <li><Link to="/web/js" activeStyle={STYLE_TAB_ACTIVE}>JS SDK</Link></li>
                    </ul>
                    {this.props.children && React.cloneElement(this.props.children, {appid: this.props.webInfo.appid})}
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


module.exports = {
    ApiTpl: ApiTpl,
    JsTpl: JsTpl,
    WebTpl: connect(mapStateToProps, mapDispatchToProps)(WebTpl),
};
