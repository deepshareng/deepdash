var React = require('react');
var connect = require('react-redux').connect;

var DSAPI = require('../lib/api.js');
var DSUTIL = require('../lib/util.js');
var DSCOMMON = require('../lib/common.js');

var Link = require('react-router').Link;
var hashHistory = require('react-router').hashHistory;

var Router = require('react-router').Router;
var Route = require('react-router').Route;
var Table = require('../components/table.jsx');

var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    // mapStateToProps must contain the full data
    return {overallInfo: state.overallInfo};
};

var mapDispatchToProps = function(dispatch) {
    var updateAbstract = function(totalPv, totalClick, totalActive, totalInstall) {
        dispatch({
            type: 'UPDATE_OVERALL_INFO',
            conf: {
                totalPv: totalPv,
                totalClick: totalClick,
                totalActive: totalActive,
                totalInstall: totalInstall,
            }, 
        });
    };
    return {
        GRAN_TEXT: {
            't': '累计所有',
            'd': '每天',
            'w': '每周',
            'm': '每月',
        },
        OS_TEXT: {
            'all': '全部设备',
            'ios': 'iOS',
            'android': "Android",
        },
        channelIndex: 0,
        sessionAppId: null,
        gotoUrl: function(url) {
            window.location.href = url;
        },
        updateGranularity: function(gran) {
            dispatch({
                type: 'UPDATE_OVERALL_INFO',
                conf: {
                    gran: gran,
                }, 
            });
        },
        updateOSPlatform: function(oSPlatform) {
            dispatch({
                type: 'UPDATE_OVERALL_INFO',
                conf: {
                    oSPlatform: oSPlatform,
                }, 
            });
        },
        updateAbstract: function(totalPv, totalClick, totalActive, totalInstall) {
            dispatch({
                type: 'UPDATE_OVERALL_INFO',
                conf: {
                    totalPv: totalPv,
                    totalClick: totalClick,
                    totalActive: totalActive,
                    totalInstall: totalInstall,
                }, 
            });
        },
        renderAbstractView: function() {
            DSAPI.getAbstract(sessionAppId, function(data) {
                updateAbstract(+data.all[0]["url-tem-overall-value"],
                  +data.all[0]["sharelink-tem-overall-value"],
                  +data.all[0]["match/open_with_params"],
                  +data.all[0]["match/install_with_params"]);
            });
        },
        updateChartData: function(data) {
            dispatch({
                type: 'UPDATE_OVERALL_INFO',
                conf: {
                  chartDataByGran: data,
                }, 
            });
        }
    };
}

var OverallTpl = React.createClass({
    getInitialState: function() {
        return {
            tableHead:[],
            tableContent: [[[]]]
        }
    },
    initChartData: function(gran) {
        var me = this;
        gran = gran || me.props.overallInfo.gran;
        DSAPI.getDist(sessionAppId, gran, function(data) {
            var x = [];
            var openlist = [];
            var installlist = [];
            var pvlist = [];
            var clicklist = [];

            var ios_openlist = [];
            var ios_installlist = [];
            var ios_pvlist = [];
            var ios_clicklist = [];

            var android_openlist = [];
            var android_installlist = [];
            var android_pvlist = [];
            var android_clicklist = [];

            for (var i in data.all) {
                x[i] = DSUTIL.addDate(i - data.all.length + 1, gran);

                openlist[i] = +data.all[data.all.length - 1 - i]["match/open_with_params"];
                installlist[i] = +data.all[data.all.length - 1 - i]["match/install_with_params"];
                pvlist[i] = +data.all[data.all.length - 1 - i]["url-tem-overall-value"];
                clicklist[i] = +data.all[data.all.length - 1 - i]["sharelink-tem-overall-value"];
            }
            for (var i in data.ios) {
                ios_openlist[i] = +data.ios[data.ios.length - 1 - i]["match/open_with_params"];
                ios_installlist[i] = +data.ios[data.ios.length - 1 - i]["match/install_with_params"];
                ios_pvlist[i] = +data.ios[data.ios.length - 1 - i]["url-tem-overall-value"];
                ios_clicklist[i] = +data.ios[data.ios.length - 1 - i]["sharelink-tem-overall-value"];
            }
            for (var i in data.android) {
                android_openlist[i] = +data.android[data.android.length - 1 - i]["match/open_with_params"];
                android_installlist[i] = +data.android[data.android.length - 1 - i]["match/install_with_params"];
                android_pvlist[i] = +data.android[data.android.length - 1 - i]["url-tem-overall-value"];
                android_clicklist[i] = +data.android[data.android.length - 1 - i]["sharelink-tem-overall-value"];
            }
            var chartDataByGran = {};
            chartDataByGran[gran] = {
                chartXAxis: x,
                chartDataActive: {
                    all: openlist, 
                    ios: ios_openlist, 
                    android: android_openlist,
                },
                chartDataInstall: {
                    all: installlist, 
                    ios: ios_installlist, 
                    android: android_installlist,
                },
                chartDataPv: {
                    all: pvlist, 
                    ios: ios_pvlist, 
                    android: android_pvlist,
                },
                chartDataClick: {
                    all: clicklist, 
                    ios: ios_clicklist, 
                    android: android_clicklist,
                },
            };

            me.props.updateChartData(chartDataByGran);
            me.renderCharts({
              gran: gran,
            });
            me.renderTable({gran:gran});
        });
    },
    renderCharts: function(options) {
        var me = this;
        options = options || {};
        var gran = options.gran || me.props.overallInfo.gran;
        var oSPlatform = options.oSPlatform || me.props.overallInfo.oSPlatform;
        var colors = [];
        var series = [];
        $(".panel.selected").each(function(){
            var $ele = $(this);
            switch($ele.attr('id')) {
                case 'panel-total-pv': 
                    series.push({
                        name: '链接展示',
                        data: me.props.overallInfo.chartDataByGran[gran].chartDataPv[oSPlatform],
                    });
                    break;
                case 'panel-total-click':
                    series.push({
                        name: '链接点击',
                        data: me.props.overallInfo.chartDataByGran[gran].chartDataClick[oSPlatform],
                    });
                    break;
                case 'panel-total-active':
                    series.push({
                        name: '累计活跃',
                        data: me.props.overallInfo.chartDataByGran[gran].chartDataActive[oSPlatform],
                    });
                    break;
                case 'panel-total-install':
                    series.push({
                        name: '累计新增',
                        data: me.props.overallInfo.chartDataByGran[gran].chartDataInstall[oSPlatform],
                    });
                    break;
            }
            colors.push($ele.find('.panel-legend-line .bottom').css('background-color'));
        });

        var rect = null;
        function drawRect(chart){
            if (rect){
               rect.element.remove();   
            }
            var xAxis = chart.xAxis[0];
            rect = chart.renderer.rect(0, chart.chartHeight - xAxis.bottom, chart.chartWidth , xAxis.bottom, 0).attr({
                'stroke-width': 0,
                stroke: '#e9f6ff',
                fill: '#e9f6ff',
                zIndex: 3,
            }).add();
        }
        $('#highcharts-area').highcharts({
            title:'',
            credits: { enabled: false, },
            chart: { type: 'spline',
                events: {
                    load: function() {
                        drawRect(this);
                    },
                    redraw: function(){
                        drawRect(this);
                    }
                },
            },
            plotOptions: {
                spline: {
                    marker: {
                        enabled: false,
                    },
                },
                series: {
                    shadow: {
                        color: '#a0a0a0', 
                        offsetY: 2
                    },
                },
            },
            legend: { enabled: false },
            colors: colors,
            xAxis: {
                categories: me.props.overallInfo.chartDataByGran[me.props.overallInfo.gran].chartXAxis,
                tickLength: 0,
                labels:{
                    style: {
                        color: '#3c99e6',
                        fontSize: '12px',
                    },
                },
            },
            yAxis: {
                title: '',
                allowDecimals: false,
                gridLineColor : '#eeeeee',
                labels:{
                    style: {
                        color: '#83caf5',
                        fontSize: '16px',
                    },
                },
            },
            series: series,
        });
    },
    renderTable: function(options){
        var me = this;
        var date = new Date();
        options = options || {};
        var gran = options.gran || me.props.overallInfo.gran;
        var oSPlatform = options.oSPlatform || me.props.overallInfo.oSPlatform;
        var head = ["日期", "链接展示", "链接点击", "累计活跃", "累计新增"];
        var pv = me.props.overallInfo.chartDataByGran[gran].chartDataPv[oSPlatform];
        var click = me.props.overallInfo.chartDataByGran[gran].chartDataClick[oSPlatform];
        var active = me.props.overallInfo.chartDataByGran[gran].chartDataActive[oSPlatform];
        var install = me.props.overallInfo.chartDataByGran[gran].chartDataInstall[oSPlatform];
        var dateList = me.props.overallInfo.chartDataByGran[me.props.overallInfo.gran].chartXAxis;
        var content = []
        for (var i = 0; i < dateList.length; i++) {
           var row = [dateList[i], pv[i], click[i], active[i], install[i]];
           content.push(row);
        }
        this.setState({
            tableHead:head,
            tableContent: content
        });
    },
    changeTimeGran: function(e) {
        // TODO: The show waiting...
        var gran = e.target.dataset.timegran;
        this.props.updateGranularity(gran); 
        this.initChartData(gran);
    },
    changeOSPlatform: function(e) {
        var oSPlatform = e.target.dataset.osplatform;
        this.props.updateOSPlatform(oSPlatform); 
        this.renderCharts({
          oSPlatform: oSPlatform,
        });
        this.renderTable({
          oSPlatform: oSPlatform,
        });
    },
    selectPanel: function(e) {
        var me = this;
        $(e.target).closest('.panel').toggleClass('selected');
        me.renderCharts();
    },
    componentWillMount: function() {
    },
    componentDidMount: function() {
        var me = this;

        //default selected panel
        if ($('.panel.selected').length <= 0) {
            $('#panel-total-pv').addClass('selected');
        }

        me.props.renderAbstractView();
        me.initChartData();

        DSGLOBAL.initComponent = function() {
            me.props.renderAbstractView();
            me.initChartData();
        };

        $("#panel-total-pv .ds-tip").popover({
            html: true, toggle: 'popover', trigger: 'hover', placement: 'top', title: '',
            content: '链接展示',
        });
        $("#panel-total-click .ds-tip").popover({
            html: true, toggle: 'popover', trigger: 'hover', placement: 'top', title: '',
            content: '链接点击',
        });
        $("#panel-total-active .ds-tip").popover({
            html: true, toggle: 'popover', trigger: 'hover', placement: 'top', title: '',
            content: '累计活跃',
        });
        $("#panel-total-install .ds-tip").popover({
            html: true, toggle: 'popover', trigger: 'hover', placement: 'top', title: '',
            content: '累计新增',
        });
    },
    render: function() {
        return (
            <div id="overall-page">
              <div className="row filter-bar">
                <div className="col-md-6 col-xs-6 text-left">
                  <div className="pull-left">
                    <div className="dropdown">
                      <a href="#" data-toggle="dropdown" aria-expanded="false" className="dropdown-toggle">{this.props.OS_TEXT[this.props.overallInfo.oSPlatform]}<span className="caret"></span></a>
                      <ul role="menu" className="dropdown-menu pull-left">
                        <li><a data-osplatform="all" href="javascript:;" onClick={this.changeOSPlatform}>{this.props.OS_TEXT['all']}</a></li>
                        <li className="divider"></li>
                        <li><a data-osplatform="android" href="javascript:;" onClick={this.changeOSPlatform}>{this.props.OS_TEXT['android']}</a></li>
                        <li><a data-osplatform="ios" href="javascript:;" onClick={this.changeOSPlatform}>{this.props.OS_TEXT['ios']}</a></li>
                      </ul>
                    </div>
                  </div>
                </div>
                <div className="col-md-6 col-xs-6">
                  <div className="pull-right">
                    <div className="dropdown">
                      <a href="#" data-toggle="dropdown" aria-expanded="false" className="dropdown-toggle">{this.props.GRAN_TEXT[this.props.overallInfo.gran]}<span className="caret"></span></a>
                      <ul role="menu" className="dropdown-menu pull-right">
                        <li><a data-timegran="d" href="javascript:;" className="time-gran" onClick={this.changeTimeGran}>{this.props.GRAN_TEXT['d']}</a></li>
                        <li><a data-timegran="w" href="javascript:;" className="time-gran" onClick={this.changeTimeGran}>{this.props.GRAN_TEXT['w']}</a></li>
                        <li className="divider"></li>
                        <li><a data-timegran="m" href="javascript:;" className="time-gran" onClick={this.changeTimeGran}>{this.props.GRAN_TEXT['m']}</a></li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
              <div className="row">
                <div className="col-md-3 col-xs-6">
                  <div className="panel" id="panel-total-pv" onClick={this.selectPanel}>
                    <div className="row">
                      <span className="col-xs-9 col-md-9">链接展示</span>
                      <span className="col-xs-3 col-md-3 text-center ds-tip"><i className="fa fa-info-circle"></i></span>
                      <span className="panel-num col-xs-12 col-md-12">{this.props.overallInfo.totalPv}</span>
                      <span className="unit text-right col-xs-12 col-md-12">/次</span>
                      <div className="panel-legend-line total-install col-xs-12 col-md-12">
                        <span className="top"></span>
                        <span className="bottom"></span>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="col-md-3 col-xs-6">
                  <div className="panel" id="panel-total-click" onClick={this.selectPanel}>
                    <div className="row">
                      <span className="col-xs-9 col-md-9">链接点击</span>
                      <span className="col-xs-3 col-md-3 text-center ds-tip"><i className="fa fa-info-circle"></i></span>
                      <span className="panel-num col-xs-12 col-md-12">{this.props.overallInfo.totalClick}</span>
                      <span className="unit text-right col-xs-12 col-md-12">/次</span>
                      <div className="panel-legend-line total-install col-xs-12 col-md-12">
                        <span className="top"></span>
                        <span className="bottom"></span>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="col-md-3 col-xs-6">
                  <div className="panel" id="panel-total-active" onClick={this.selectPanel}>
                    <div className="row">
                      <span className="col-xs-9 col-md-9">累计活跃</span>
                      <span className="col-xs-3 col-md-3 text-center ds-tip"><i className="fa fa-info-circle"></i></span>
                      <span className="panel-num col-xs-12 col-md-12">{this.props.overallInfo.totalActive}</span>
                      <span className="unit text-right col-xs-12 col-md-12">/人</span>
                      <div className="panel-legend-line total-install col-xs-12 col-md-12">
                        <span className="top"></span>
                        <span className="bottom"></span>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="col-md-3 col-xs-6">
                  <div className="panel" id="panel-total-install" onClick={this.selectPanel}>
                    <div className="row">
                      <span className="col-xs-9 col-md-9">累计新增</span>
                      <span className="col-xs-3 col-md-3 text-center ds-tip"><i className="fa fa-info-circle"></i></span>
                      <span className="panel-num col-xs-12 col-md-12">{this.props.overallInfo.totalInstall}</span>
                      <span className="unit text-right col-xs-12 col-md-12">/人</span>
                      <div className="panel-legend-line total-install col-xs-12 col-md-12">
                        <span className="top"></span>
                        <span className="bottom"></span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="row">
                <div className="col-md-12 col-xs-12">
                  <div id="user-new-chart" className="panel panel-default">
                    <div id="highcharts-area">
                      <div className="text-center loader-inner ball-pulse">
                        <div></div>
                        <div></div>
                        <div></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="row">
                <div className="col-md-12 col-xs-12">
                    <div className="panel panel-default">
                        <Table.regular id="data-table" head={this.state.tableHead} content={this.state.tableContent}/>                     
                    </div>
                </div>
              </div>
            </div>);
    },
});

module.exports = {
    OverallTpl: connect(mapStateToProps, mapDispatchToProps)(OverallTpl),
};
