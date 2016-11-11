var React = require('react');
var commonStyle = require('../styles/common.js')

var Panel = {};


Panel.regular = React.createClass({
  render: function() {
    var classes = "panel panel-default " + this.props.theme;
    return (  
      <div className={classes}>
        <div className="panel-heading">
          <h3 className="panel-title">{this.props.title}</h3>
        </div>
        <div className="panel-body">
          <dl className="dl-horizontal h5">
            {this.props.content.map((item) => 
                [
                <dt className="text-muted">{item.title}</dt>,
                <dd id={item.id} className="text-primary">{item.content}</dd>
                ]
            )}
          </dl>
        </div>
      </div>
    );
  }
});

Panel.row = React.createClass({
  render: function() {
    return (
      <div className="row">
        <div className={this.props.size}>
          <Panel.regular title={this.props.title} theme={this.props.theme} content={this.props.content}/>
        </div>
      </div>
    )
  }
});

module.exports = Panel;
