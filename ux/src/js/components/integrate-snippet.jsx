var React = require('react');
var ReactDOM = require('react-dom');
var DSCOMMON = require('../lib/common.js');

var STYLE_STEPS_HEADER = require('../styles/common.js').STYLE_STEPS_HEADER;

var HorizontalProgressBar = React.createClass({
    render: function() {
        var datas = [
            {name: '创建', path: '/#/create-app'}, 
            {name: '配置', path: '/#/configure'}, 
            {name: 'APP集成', path: '/#/app'}, 
            {name: '网页集成', path: '/#/web'}, 
            {name: '完成', path: '/#/complete'},
        ];
        var stepIndex = this.props.stepIndex;
        var headers = datas.map(function(d, i) {
            var btnClassName = 'btn btn-default btn-circle ds-integrate-status-icon';
            var blockClassName = "ds-integrate-steps-header inactive";
            if (stepIndex > i) {
                btnClassName = 'btn btn-success btn-circle ds-integrate-status-icon';
                blockClassName = "ds-integrate-steps-header";
            }
            return (
                <a href={sessionAppId==''?'/#/create-app':datas[i].path}>
                    <div key={i} className={blockClassName}>
                        <span className="ds-integrate-steps-text">{datas[i].name}</span>
                        <span className={btnClassName}>
                            {i + 1} 
                        </span>
                    </div>
                </a>
                ); 
        });
        return (
            <div className="ds-integrate-steps-header-container">
                <div className="ds-integrate-steps-header-container-fixed">
                    <div>{headers}</div>
                </div>
            </div>
        ); 
    },
});

var CompleteCreateTab = React.createClass({
    render: function() {
        return (
          <div className="row">
            <div className="col-lg-8 col-lg-offset-2">
              <div className="panel panel-default">
                <div className="panel-heading">
                  <h3 className="panel-title"> 
                    <span className="btn btn-success btn-circle ds-integrate-status-icon">
                        <i className="fa fa-check"></i>
                    </span>
                    <span>1.你已经创建了APP -  
                    <span className="show-appname">{this.props.appName}</span></span>
                  </h3>
                </div>
              </div>
            </div>
          </div>
        );        
    },
});

var CompleteConfTab = React.createClass({
    render: function() {
        return (
          <div className="row">
            <div className="col-lg-8 col-lg-offset-2">
              <div className="panel panel-default">
                <div className="panel-heading">
                  <h3 className="panel-title">
                     <span className="btn btn-success btn-circle ds-integrate-status-icon">
                         <i className="fa fa-check"></i>
                     </span>
                     <span>2.你已配置了APP信息</span>
                  </h3>
                </div>
              </div>
            </div>
          </div>
        );        
    },
});

var CompleteAppTab = React.createClass({
    render: function() {
        return (
          <div className="row">
            <div className="col-lg-8 col-lg-offset-2">
              <div className="panel panel-default">
                <div className="panel-heading">
                  <h3 className="panel-title"> <span className="btn btn-success btn-circle ds-integrate-status-icon"><i className="fa fa-check"></i></span><span>3.你已集成了SDK</span></h3>
                </div>
              </div>
            </div>
          </div>
        );        
    },
});


var CopyCodeBtn = React.createClass({
    getInitialState: function() {
        return {eId: 'id'+ Math.round(Math.random()*10000000)};
    },
    componentDidMount: function() {
        //var cb = new Clipboard('#''.btn-copy-code', {
        var cb = new Clipboard('#'+this.state.eId, {
            target: function(trigger) {
                return trigger.parentNode; 
            }, 
        });
        cb.on('success', function(e) {
            e.clearSelection(); 
            DSCOMMON.alert('代码已复制到剪贴板');
        });
    },
    render: function() {
        return (
            <button id={this.state.eId} className="btn btn-outline btn-default btn-sm btn-copy-code" type="button">
            复制
            </button>
        );        

        //<i className="fa fa-clipboard"></i>
    },
});

module.exports = {
    HorizontalProgressBar: HorizontalProgressBar,
    CompleteCreateTab: CompleteCreateTab,
    CompleteConfTab: CompleteConfTab,
    CompleteAppTab: CompleteAppTab,
    CopyCodeBtn: CopyCodeBtn,
};
