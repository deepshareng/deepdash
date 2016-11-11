var React = require('react');
var connect = require('react-redux').connect;

var DSCOMMON = require('../lib/common.js');
var DSAPI = require('../lib/api.js');

var getDateStr = function(daysOffset) {
    var now = new Date();
    var before = new Date(now.getTime() - daysOffset * 24 * 60 * 60 * 1000);

    return before.toJSON().slice(0, 10);
};

var defaultHead = [
    '渠道类型', '渠道名称', '渠道编号及备注', '点击量',
    '新用户下载', '老用户打开',
    // '三日留存', '七日留存',
];
var defaultMetric = [
    'typename', 'channelname', 'remark', 'sharelink-tem-overall-value',
    'match/install_with_params', 'match/open_with_params',
    //'3-day-retention', '7-day-retention',
];
var defaultMetricAppend = [
    '', '', '', '',
    '', '',
    //'%', '%',
];
var defaultEvents = [
    'sharelink-tem-overall-value',
    'match/install_with_params',
    'match/open_with_params',
    //'3-day-retention', '7-day-retention',
];

var parseCheckedMetricsDisplay = function() {
    var v = $('#metrics input[type=checkbox]:checked');
    var metrics = [];
    for (var i = 0; i < v.length; i++) {
        metrics.push(v[i].dataset.display);
    }

    return metrics;
};
var getDateStart = function() {
    return $('#datepicker-start').val();
};
var getDateEnd = function() {
    return $('#datepicker-end').val();
};
var parseCheckedMetrics = function() {
    var v = $('#metrics input[type=checkbox]:checked');
    var metrics = [];
    for (var i = 0; i < v.length; i++) {
        metrics.push(v[i].value);
    }

    return metrics;
};

var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    // mapStateToProps must contain the full data
    return {marketingInfo: state.marketingInfo};
};

var mapDispatchToProps = function(dispatch) {
    var updateMetrics = function(metrics) {
        dispatch({
            type: 'UPDATE_MARKETING_INFO',
            conf: {
                metrics: metrics,
            }, 
        });
    };
    var updateTable = function(tableHeads, tableBodyKeys, tableData){
        dispatch({
            type: 'UPDATE_MARKETING_INFO',
            conf: {
                tableHeads: tableHeads,
                tableBodyKeys: tableBodyKeys,
                tableData: tableData,
            }, 
        });
    };
    var updateDateStr = function(dateStartStr, dateEndStr) {
        dispatch({
            type: 'UPDATE_MARKETING_INFO',
            conf: {
                dateStartStr: dateStartStr,
                dateEndStr: dateEndStr,
            }, 
        });
    };
    var renderStatisticsTable = function() {
        var queryData = {
            event: defaultEvents.concat(parseCheckedMetrics()).join(','),
            start: getDateStart(),
            end: getDateEnd(),
        };

        DSAPI.getChannelStatistics(sessionAppId, queryData, function() {
            $('.loader-inner').removeClass('hide');
        }, function(data) {
            var tableData = data.data;

            var tableHeads = defaultHead.concat(parseCheckedMetricsDisplay());
            var tableBodyKeys = defaultMetric.concat(parseCheckedMetrics());

            if (tableData.length <= 0) {
                $('#default-blank-channel-statistics').removeClass('hide');
            } else {
                $('#default-blank-channel-statistics').addClass('hide');
            }
            updateTable(tableHeads, tableBodyKeys, tableData);
        }, function() {
            DSCOMMON.alert('获取数据失败，请刷新页面');
        }, function() {
            $('.loader-inner').addClass('hide');
        });
    };
    var fetchMetrics = function() {
        // TODO: call the api to reder, metrics dom area!!
        // it is ok that the fetchmetric and fetch statistics 'not sync'!!!,
        // it can work right without metrics from server!!!
        //
        DSAPI.getEvents(sessionAppId, function(data) {
            var events = data.eventlist;
            var metrics = [];

            for (var i = 0; i < events.length; i++) {
                // filter only with counters: customized attribute
                if (/^\/v2\/counters\//.test(events[i]['event'])) {
                    metrics.push(events[i]);
                } 
            }
            updateMetrics(metrics);
        }, function() {});
    };
    return {
        initializeConf: function() {
            $('#datepicker').datepicker({
                format: 'yyyy-mm-dd', 
            });

            $('#time-span-quick-picker').on('change', function() {
                var gran = $(this).val();
                var spanConf = {
                    '0': [0, 0], 
                    '1': [1, 1],
                    '7': [6, 0],
                    '30': [29, 0],
                    '180': [179, 0],
                    '365': [364, 0],
                };

                var start = '';
                var end = '';
                if (gran in spanConf) {
                    start = getDateStr(spanConf[gran][0]);
                    end = getDateStr(spanConf[gran][1]);
                } else if (gran === 'all') {
                    end = getDateStr(0);
                } else {
                    start = getDateStr(7);
                    end = getDateStr(1);
                }
                updateDateStr(start, end);
            });

            // set default date of datepicker
            updateDateStr(getDateStr(6), getDateStr(0));
            $('#time-span-quick-picker').val(7);

            renderStatisticsTable();
            fetchMetrics();
        },
        renderStatisticsTableAdapter: function() {
            renderStatisticsTable();
        },
        updateDateStartStr: function(e) {
            dispatch({
                type: 'UPDATE_MARKETING_INFO',
                conf: {
                    dateStartStr: e.target.value,
                }, 
            });
        },
        updateDateEndStr: function(e) {
            dispatch({
                type: 'UPDATE_MARKETING_INFO',
                conf: {
                    dateEndStr: e.target.value,
                }, 
            });
        },
    };
}

var MarketingTpl = React.createClass({
    
    componentWillMount: function() {
    },
    queryMarketingData: function() {
        this.props.renderStatisticsTableAdapter();
    },
    componentDidMount: function() {
        this.props.initializeConf();
        DSGLOBAL.initComponent = this.props.initializeConf;
    },
    render: function() {
        //Metrics
        var self = this;
        var metricsTpl;
        if (this.props.marketingInfo.metrics.length > 0) {
            metricsTpl = this.props.marketingInfo.metrics.map(function(data) {
                return <label className="checkbox-inline">
                            <input type="checkbox" data-display={data['display']} value={data['event']}/>
                            {data['display']}
                        </label>;
            })
        } else {
            metricsTpl = <span className="text-muted">还未添加指标</span>; 
        }

        return (
            <div>
                <div className="row">
                  <div className="col-lg-12">
                    <h1 className="page-header">渠道检测</h1>
                  </div>
                </div>
                <div className="row col-lg-12 ds-marketing-selector-area">
                  <div className="col-lg-12 ds-marketing-selector">
                    <div className="selector-label">时间区间:</div>
                    <div id="datepicker" className="input-daterange input-group selector-item">
                      <input id="datepicker-start" value={this.props.marketingInfo.dateStartStr} onChange={this.props.updateDateStartStr} type="text" size="6" name="start" className="input-sm form-control"/><span className="input-group-addon">-&gt;</span>
                      <input id="datepicker-end" value={this.props.marketingInfo.dateEndStr} onChange={this.props.updateDateEndStr} type="text" size="6" name="end" className="input-sm form-control"/>
                    </div>
                    <div className="selector-item">
                      <select id="time-span-quick-picker" className="form-control input-sm">
                        <option value="0">今天</option>
                        <option value="1">昨天</option>
                        <option value="7">最近7天</option>
                        <option value="30">最近30天</option>
                        <option value="180">最近180天</option>
                        <option value="365">最近1年</option>
                        <option value="all">所有数据</option>
                      </select>
                    </div>
                  </div>
                  <div className="col-lg-12 ds-marketing-selector">
                    <div className="selector-label">关键指标:</div>
                    <div id="metrics" className="selector-item">{metricsTpl}</div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-lg-12">
                    <ul className="nav nav-tabs">
                      <li className="active"><a href="#">统计数据</a></li><a onClick={this.queryMarketingData} className="btn btn-primary btn-outline pull-right">查询 </a>
                    </ul>
                    <div className="panel panel-default">
                      <div className="panel-body">
                        <div className="text-center loader-inner ball-pulse">
                          <div></div>
                          <div></div>
                          <div></div>
                        </div>
                        <div className="table-responsive">
                          <table id="channel-statistics-table" className="table table-striped">
                            <thead key="thead"><tr>
                            {
                                this.props.marketingInfo.tableHeads.map(function (tableHead, i) {
                                    return <th key={i}>{tableHead}</th>;
                                })
                            }
                            </tr>
                            </thead>
                            <tbody key="tbody">
                            {
                                this.props.marketingInfo.tableData.map(function (data, i) {
                                    var nodes = self.props.marketingInfo.tableBodyKeys.map(function(key, j) {
                                            var v = data[key] + (defaultMetricAppend[j] || '');
                                            return <td key={j}>{v}</td>;
                                        });
                                    return <tr key={i}>{nodes}</tr>;
                                })
                            }
                            </tbody>
                          </table>
                        </div>
                        <div id="default-blank-channel-statistics" className="text-center default-blank hide">当前应用还没有渠道统计数据 ^_^</div>
                      </div>
                    </div>
                  </div>
                </div>
            </div>);
    },
});

module.exports = {
    MarketingTpl: connect(mapStateToProps, mapDispatchToProps)(MarketingTpl),
};
