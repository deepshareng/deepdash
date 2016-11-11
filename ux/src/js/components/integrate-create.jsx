var React = require('react');
var ReactDOM = require('react-dom');
var connect = require('react-redux').connect;

var DSAPI = require('../lib/api.js');
var DSCOMMON = require('../lib/common.js');

var HorizontalProgressBar = require('./integrate-snippet.jsx').HorizontalProgressBar;


var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    return {
        createInfo: state.createInfo, 
    };
};

var mapDispatchToProps = function(dispatch) {
    return {
        onInputChange: function(e) {
            dispatch({
                type: 'UPDATE_APP_NAME',
                appname: e.target.value,
            });
        },
        initializeCreateApp: function() {
        },
    };
};

var CreateAppTpl = React.createClass({
    getInitialState: function() {
        return {};                
    },
    componentDidMount: function() {
        var self = this;
        
        if (sessionAppId) {
            DSAPI.getAppInfo(sessionAppId, function(data) {
                self.setState({
                    sessionAppId: sessionAppId,
                    appName: data.appName});
            });
        }

        DSGLOBAL.initComponent = function() {
            if (sessionAppId) {
                DSAPI.getAppInfo(sessionAppId, function(data) {
                    self.setState({
                        sessionAppId: sessionAppId,
                        appName: data.appName});
                });
            }
        };
    },
    onClickNext: function(e) {
        if (this.props.createInfo.inputAppName.length <= 0) {
            $('#appname-input-group').addClass('has-error');
            return;
        }

        DSAPI.createApp({appName: this.props.createInfo.inputAppName}, function(data) {
            if (data.appid) {
                // update global current appid
                sessionAppId = data.appid; 
                DSCOMMON.getApps(function(data) {
                    appNames = data.appNames;
                    DSCOMMON.selectApp(data.appIndex);
                });
                location.href = '#/configure';
            } else {
                // TODO: handle error
            }
        }, function(xhr) {
            if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                DSCOMMON.showLoginModal()
            }
        });  
    },
    onClickDirect: function() {
        location.href = '#/configure';
    },
    onInputAppName: function(e) {
        if (e.which === 13) {
            this.onClickNext();
            return false;
        } 
    },
    render: function() {
        return (
        <div>
            <HorizontalProgressBar stepIndex={1}/> 
            <div className="row ds-integrate-steps-content">
                <div className="col-lg-8 col-lg-offset-2">
                    <div className="panel panel-default">
                        <div className="panel-heading">
                            <h3 className="panel-title">
                                <span className="btn btn-info btn-circle ds-integrate-status-icon">
                                    <i className="fa fa-arrow-right"></i>
                                </span>
                                <span> 1.创建APP </span>
                            </h3>
                        </div>
                        <div className="panel-body">
                            <form role="form">
                                <div id="appname-input-group" className="form-group">
                                    <label className="control-label" htmlFor="appname">应用名称</label>
                                    <input className="form-control" type="text" value={this.props.createInfo.inputAppName || ''} onKeyDown={this.onInputAppName} onChange={this.props.onInputChange}/>
                                    <p className="text-info help-block">请输入应用名称</p>
                                </div>
                                <div>
                                    <a id="addapp-next-step" onClick={this.onClickNext} className="btn btn-primary">下一步</a>
                                </div>
                                <div id="direct-next-area" className={this.state.sessionAppId ? '' : 'hide'}>
                                    <hr/>
                                    <div className="form-group">
                                        您已经添加过应用，想集成当前应用 
                                        <strong className="text-info" style={{marginLeft: '0.8rem', marginRight: '0.8rem'}}>[{this.state.appName}]</strong>
                                        ，请点击‘直接集成’
                                    </div>
                                    <div>
                                        <a id="direct-next-step" onClick={this.onClickDirect} className="btn btn-primary">直接集成</a>
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>    
        </div>
        );        
    },
});

module.exports = {
    CreateAppTpl: connect(mapStateToProps, mapDispatchToProps)(CreateAppTpl)
};
