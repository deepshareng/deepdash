var React = require('react');
var Table = {};


Table.regular = React.createClass({
  render: function() {
    return (
      <div className="table-responsive">
        <table className="table table-striped table-hover table-bordered border-radius">
          <thead>
          <tr>
            {this.props.head.map((item) => 
              <th>{item}</th>
            )}
          </tr>
          </thead>
          <tbody>
          {this.props.content.map((row) => 
            <tr>
              {row.map((col) => 
                <td>{col}</td>
              )}
            </tr>
          )}
          </tbody>
        </table>
      </div>
    );
  }
});



module.exports = Table;
