var React = require('react');
var ReactDOM = require('react-dom');
var connect = require('react-redux').connect;

var DSCOMMON = require('../lib/common.js');

var simulatorConf = require('../reducers/simulator.js').simulatorConf;

var mapStateToProps = function(state={}, ownProps) {
    return {simulatorInfo: state.simulatorInfo};
};

var mapDispatchToProps = function(dispatch) {
    return {
        onChangePlatform: function(e) {
            dispatch({
                type: 'UPDATE_SIMULATOR_TYPE',
                platform: e.target.dataset.platform,
            });
        },
        onChangeDownload: function(e) {
            dispatch({
                type: 'UPDATE_SIMULATOR_TYPE',
                download: e.target.dataset.download,
            });
        },
        onChangeInstall: function(e) {
            dispatch({
                type: 'UPDATE_SIMULATOR_TYPE',
                install: e.target.dataset.install,
            });
        },
        initializeCreateApp: function() {
        },
    };
};

var SimulatorTpl = React.createClass({
    getInitialState: function() {
        return {};                
    },
    componentDidMount: function() {
    },
    render: function() {
        return (
        <div>
            <div className="row">
                <div className="col-sm-5 col-sm-offset-1 simulator-nav">
                    <span className="title"><i className="fa fa-circle-thin" aria-hidden="true"></i>操作系统</span>
                    <span onClick={this.props.onChangePlatform} data-platform="IOS9" className={this.props.simulatorInfo.platform === 'IOS9'?"btn-radio selected" : "btn-radio"}>iOS 9 及以上</span>
                    <span onClick={this.props.onChangePlatform} data-platform="IOS8" className={this.props.simulatorInfo.platform === 'IOS8'?"btn-radio selected" : "btn-radio"}>iOS 8 及以下</span>
                    <span onClick={this.props.onChangePlatform} data-platform="ANDROID" className={this.props.simulatorInfo.platform === 'ANDROID'?"btn-radio selected" : "btn-radio"}>Android</span>
                    <span className="title"><i className="fa fa-circle-thin" aria-hidden="true"></i>是否开启应用宝</span>
                    <span onClick={this.props.onChangeDownload} data-download="YYB" className={this.props.simulatorInfo.download === 'YYB'?"btn-radio selected" : "btn-radio"}>是</span>
                    <span onClick={this.props.onChangeDownload} data-download="DIRECT" className={this.props.simulatorInfo.download === 'DIRECT'?"btn-radio selected" : "btn-radio"}>否</span>
                    <span className="title"><i className="fa fa-circle-thin" aria-hidden="true"></i>是否安装APP</span>
                    <span onClick={this.props.onChangeInstall} data-install="INSTALLED" className={this.props.simulatorInfo.install === 'INSTALLED'?"btn-radio selected" : "btn-radio"}>是</span>
                    <span onClick={this.props.onChangeInstall} data-install="NEW" className={this.props.simulatorInfo.install === 'NEW'?"btn-radio selected" : "btn-radio"}>否</span>
                    <div style={{clear: 'both'}}></div>
                    <div className="tips">
                        <span>注意事项：</span><br/>
                        <span>以上为基于微信平台的跳转流程，其他平台跳转流程有少许差异，请用户自行尝试。</span>
                    </div>
                </div>
                <div className="col-sm-6">
                    <div className="panel panel-default no-border panel-simulator">
                        <div className="panel-body">
                            <iframe src={this.props.simulatorInfo.url} width="472" height="922" allowTransparency="true" frameborder="0" style={{border: 0}}></iframe>
                        </div>
                    </div>
                </div>
            </div>    
        </div>
        );        
    },
});

module.exports = {
    SimulatorTpl: connect(mapStateToProps, mapDispatchToProps)(SimulatorTpl)
};
