var React = require('react');
var ReactDOM = require('react-dom');

var DSAPI = require('../lib/api.js');
var HorizontalProgressBar = require('./integrate-snippet.jsx').HorizontalProgressBar;

var CompleteTpl = React.createClass({
    onClickPrev: function() {
        location.href = '#/web';
    },
    componentDidMount: function() {
        var self = this;
        DSAPI.getAppInfo(sessionAppId, function(data) {
            self.setState({
                appName: data.appName});
        });

        DSGLOBAL.initComponent = function() {
            DSAPI.getAppInfo(sessionAppId, function(data) {
                self.setState({
                    appName: data.appName});
            });
        };
    },
    render: function() {
        return (
          <div>
            <HorizontalProgressBar stepIndex={5}/>
            
            <div className="row ds-integrate-steps-content">
              <div className="col-lg-8 col-lg-offset-2">
                <div className="panel panel-default">
                  <div className="panel-body text-center ds-integrate-final">
                      您已完成DeepShare的集成流程，如有问题，请联系
                      <span className="text-info">support@deepshare.io</span>
                  </div>
                  <div className="panel-footer">
                        <div><a onClick={this.onClickPrev} className="btn btn-primary">上一步</a></div>
                  </div>
                </div>
              </div>
            </div>
          </div>   
        );        
    },

});


module.exports = {
    CompleteTpl: CompleteTpl,
};
