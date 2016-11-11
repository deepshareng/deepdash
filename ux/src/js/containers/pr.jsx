var DSCOMMON = require('../lib/common.js');
var DSAPI = require('../lib/api.js');

var React = require('react');
var connect = require('react-redux').connect;

var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    return {
        prInfo: state.prInfo, 
    };
};

var mapDispatchToProps = function(dispatch) {
    return {
        initializePrInfo: function() {
            DSAPI.getChannelList(sessionAppId, function(data) {
                dispatch({
                    type: 'UPDATE_PR_INFO',
                    prInfo: data.data,
                });
            }, function() {
                //displayNoUrlData();
            });
        },
        onClickSort: function(e) {
            var sortType = e.target.id.split('_')[1];              
            dispatch({
                type: 'UPDATE_PR_ORDER',
                sortType: sortType,
            });
        },
    };
};

var PrTpl = React.createClass({
    getInitialState: function() {
        return {};                
    },
    componentDidMount: function() {
        this.props.initializePrInfo();
        DSGLOBAL.initComponent = this.props.initializePrInfo;
    },
    onInputAppName: function(e) {
        if (e.which === 13) {
            this.onClickNext();
            return false;
        } 
    },
    onClickDeleteChannel: function(e) {
        var self = this;
        var channelName = e.target.dataset.channelName;
        var confirmDelete = confirm('确认要删除推广链接: ' + channelName + ' 吗？');

        if (!confirmDelete) {
            return;
        }

        DSAPI.deleteChannel(sessionAppId, channelName, function() {
            DSCOMMON.showModal('删除推广链接', '删除推广链接成功!');
        }, function(xhr) {
            if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                DSCOMMON.showModal('删除推广链接', '没有权限，请确认账号权限或重新登录');
            } else {
                DSCOMMON.showModal('删除推广链接', '删除推广链接失败<br/>请确认网络正常，再重新尝试');
            }
        }, function() {
            self.props.initializePrInfo();
        });
    },
    onClickShowQRCode: function(e) {
        // clear -> loading -> show
        // show is handled by modal
        // FIXME: write in react style
        $('#pr-link-qrcode').html('').qrcode(e.target.dataset.prLink);
    },
    onClickDownloadQRCode: function(e) {
        var qrCanvas = $('#pr-link-qrcode canvas');  
        if (qrCanvas.length < 1) {
            DSCOMMON.alert('不能下载二维码，请重新查看二维码或刷新页面');
        } else {
            var qrfile = qrCanvas[0].toDataURL("image/png")
                .replace(/^data:image\/[^;]*/, 'data:application/octet-stream');
            e.target.download = 'link-qrcode.png';
            e.target.href =  qrfile;
        }
    },
    onClickShowAddChannel: function(e) {
        var self = this;
        DSAPI.getAppInfo(sessionAppId, function(data) {
            self.state.download_title = data.download_title;
            self.state.download_msg = data.download_msg;
        });
    },
    onClickAddChannel: function(e) {
        var self = this;
        if (self.state.allowAdd === false) {
            return;
        } else {
            self.state.allowAdd = false;
        }

        // TODO: check input
        var channelType = $('#channelType').val();
        var channelName = $('#channelName').val();
        var channelLabel = $('#channelLabel').val();

        if (channelType === '' || channelName === '') {
            DSCOMMON.alert('渠道名称、渠道类型都不能为空');
            self.state.allowAdd = true;
            return;
        }

        // Check inappdata not added
        if ($('#inapp-data-key').val() || $('#inapp-data-val').val()) {
            DSCOMMON.alert('[应用内信息]中有输入了但未添加的字段，<br/>如果想使其生效请点击\'添加\'按钮，或者将输入框清空');
            self.state.allowAdd = true;
            return;
        }

        var channels = [channelType + '_' +
                        channelName + '_' +
                        channelLabel];

        var getInappData = function() {
            var dom = $('#inapp-data-list td');
            var inappDataObj = {};
            // ignore the '删除' 按钮
            for (var i=0; i < dom.length; i = i + 3) {
                inappDataObj[$(dom[i]).text()] = $(dom[i + 1]).text(); 
            }
            return inappDataObj;
        };
        var clearAddChannelForm = function() {
            $('#channelType').val('');
            $('#channelName').val('');
            $('#channelLabel').val('');

            $('#inapp-data-list').html('');
        };

        var param = {
            sender_id: '',
            is_permanent: true,
            use_shortid: $('#use-shortid')[0].checked,
            channels: channels,
            download_title: self.state.download_title,
            download_msg: self.state.download_msg,
            inapp_data: getInappData(),
            download_url_ios: $('#downloadUrlIos').val(),
            download_url_android: $('#downloadUrlAndroid').val(),
        };

        DSAPI.generateChannelUrl(sessionAppId, param, function(data) {
            // alert success and insert 
            DSAPI.insertChannel(sessionAppId, {
                'channelname': channels[0],
                'channelurl': data.url,
            }, function(result) {
                if (!result.error) {
                    DSCOMMON.alert('插入推广链接成功!'); 
                    self.props.initializePrInfo();
                    $('#myModal').modal('hide');
                    clearAddChannelForm();
                } else {
                    DSCOMMON.showModal('添加推广链接', result.error);
                }
            }, function(xhr) {
                if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                    DSCOMMON.alert('没有权限，请确认账号权限或重新登录');
                } else {
                    DSCOMMON.alert('插入推广链接失败，请重试'); 
                }
            }, function() {
                self.state.allowAdd = true;
            });
        }, function() {
            DSCOMMON.alert('创建推广链接失败，请重试'); 
            self.state.allowAdd = true;
        });

    },
    onClickAddInappData: function(e) {
        var dataKey = $('#inapp-data-key').val();
        var dataVal = $('#inapp-data-val').val(); 

        // TODO: check value
        if (!dataKey) {
            DSCOMMON.alert('请不要填写空的key值');
            return;
        }

        $('#inapp-data-list').append(
            '<tr>' +
            '<td>' + dataKey + '</td>' +
            '<td>' + dataVal + '</td>' +
            '<td><a class="link btn-delete-inapp-data" href="javascript:;">删除</a></td>' +
            '</tr>');

        $('#inapp-data-key').val('');
        $('#inapp-data-val').val('');

        $('.btn-delete-inapp-data').on('click', function() {
            // !! the <tr> that the <a> in <td> in <tr>
            $(this).parent().parent().remove(); 
        });
    },
    render: function() {
        var self = this;
        var nodes = this.props.prInfo.links.map(function(link, i) {
            var names = link.channelname.split('_');
            var cType = names[0];
            var cName = names[1];
            var cAppd = names[2];
            return (
                <tr key={link.channelurl}>
                <td>{ i+1 }</td>
                <td>{ cType }</td>
                <td>{ cName }</td>
                <td>{ cAppd }</td>
                <td>{ link.channelurl }</td>
                <td><a className="link btn-delete-channel" href="javascript:;" data-channel-name={link.channelname} onClick={self.onClickDeleteChannel}>删除</a>
                <a className="link btn-show-qrcode" href="javascript:;" data-toggle="modal" data-target="#modal-qrcode" data-pr-link={link.channelurl} onClick={self.onClickShowQRCode}>查看二维码</a></td>
                </tr>
            );
        }); 

        if (this.props.prInfo.links.length <= 0) {
            $('#default-blank-channel-info').removeClass('hide');
        } else {
            $('#default-blank-channel-info').addClass('hide');
        }

        return (
        <div>
            <div className="row">
              <div className="col-lg-12">
                <h1 className="page-header">推广链接</h1>
              </div>
            </div>
            <div className="row">
              <div className="col-lg-12">
                <ul className="nav nav-tabs">
                  <li className="active"><a href="javascript:;">链接信息</a></li>
                  <div className="pull-right">
                    <div className="btn-toolbar">
                      <div className="btn-group">
                        <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" className="btn btn-primary btn-outline dropdown-toggle">排序 <span className="caret"></span></button>
                        <ul className="dropdown-menu">
                          <li><a id="sort_time" className="sort" onClick={this.props.onClickSort}>创建时间</a></li>
                          <li><a id="sort_name" className="sort" onClick={this.props.onClickSort}>渠道名称</a></li>
                          <li><a id="sort_type" className="sort" onClick={this.props.onClickSort}>渠道类型</a></li>
                        </ul>
                      </div>
                      <div className="btn-group">
                          <a id="btn-show-add-channel" data-toggle="modal" data-target="#myModal" className="btn btn-primary btn-outline" onClick={this.onClickShowAddChannel}>添加推广链接</a>
                      </div>
                    </div>
                  </div>
                  <div id="myModal" tabIndex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true" style={{display: 'none'}} className="modal fade">
                    <div className="modal-dialog">
                      <div className="modal-content">
                        <div className="modal-header">
                          <button type="button" data-dismiss="modal" aria-hidden="true" className="close">×</button>
                          <h4 id="myModalLabel" className="modal-title">添加推广链接</h4>
                        </div>
                        <div className="modal-body">
                          <form role="form">
                            <div className="form-group">
                              <label htmlFor="channelType" className="control-label">渠道类型</label>
                              <input id="channelType" type="text" className="form-control"/>
                              <p className="text-info help-block">如: online, offline等</p>
                            </div>
                            <div className="form-group">
                              <label htmlFor="channelName" className="control-label">渠道名称</label>
                              <input id="channelName" type="text" className="form-control"/>
                              <p className="text-danger help-block">如: weixin, baidu, wap等</p>
                            </div>
                            <div className="form-group">
                              <label htmlFor="channelLabel" className="control-label">渠道编号及备注</label>
                              <input id="channelLabel" type="text" className="form-control"/>
                            </div>
                            <div className="form-group">
                              <label htmlFor="downloadUrlAndroid" className="control-label">Android下载地址</label>
                              <input id="downloadUrlAndroid" type="text" className="form-control"/>
                              <p className="text-info help-block">可选，为此推广链接自定义下载包地址</p>
                            </div>
                            <div className="form-group">
                              <label htmlFor="downloadUrlIos" className="control-label">iOS下载地址</label>
                              <input id="downloadUrlIos" type="text" className="form-control"/>
                              <p className="text-info help-block">可选，为此推广链接自定义下载包地址</p>
                            </div>
                            <div className="checkbox">
                              <label className="control-label">
                                <input id="use-shortid" name="use-shortid" type="checkbox"/>生成短连接
                              </label>
                              <p className="text-info help-block">生成的链接更短，适用于短信、微博推广</p>
                            </div>
                            <div className="form-group">
                              <label htmlFor="channelInAppData" className="control-label">应用内信息</label>
                              <p className="text-danger help-block">(可选)</p>
                              <div className="table-responsive">
                                <table className="inapp-data-table table">
                                  <thead>
                                    <tr>
                                      <th>key</th>
                                      <th>value</th>
                                      <th>&nbsp;&nbsp;&nbsp;&nbsp; </th>
                                    </tr>
                                  </thead>
                                  <tbody id="inapp-data-list"></tbody>
                                </table>
                              </div>
                              <div className="row">
                                <div className="col-sm-4 col-xs-4">
                                  <input id="inapp-data-key" type="text" className="form-control"/>
                                </div>
                                <div className="col-sm-4 col-xs-4">
                                  <input id="inapp-data-val" type="text" className="form-control"/>
                                </div>
                                <div className="col-sm-4 col-xs-4 text-center">
                                  <button id="btn-add-inapp-data" type="button" className="btn btn-default" onClick={this.onClickAddInappData}>添加</button>
                                </div>
                                <div style={{clear: 'both'}}></div>
                              </div>
                            </div>
                          </form>
                        </div>
                        <div className="modal-footer">
                          <button type="button" data-dismiss="modal" className="btn btn-default">关闭</button>
                          <button id="btn-add-channel" type="button" className="btn btn-primary" onClick={this.onClickAddChannel}>提交</button>
                        </div>
                      </div>
                    </div>
                  </div>
                </ul>
                <div className="panel panel-default">
                  <div className="panel-body">
                    <div className="table-responsive">
                      <table id="channel-info-table" className="table table-striped">
                        <thead>
                          <tr>
                            <th>#</th>
                            <th>渠道类型</th>
                            <th>渠道名称</th>
                            <th>渠道编号及备注</th>
                            <th>渠道url </th>
                            <th>操作</th>
                          </tr>
                        </thead>
                        <tbody>{ nodes }</tbody>
                      </table>
                    </div>
                    <div id="default-blank-channel-info" className="text-center default-blank hide">当前应用还没有添加推广链接 ^_^</div>
                  </div>
                  <div id="modal-qrcode" tabIndex="-1" role="dialog" aria-labelledby="modal-qrcode-label" aria-hidden="true" style={{display: 'none'}} className="modal fade">
                    <div className="modal-dialog">
                      <div className="modal-content">
                        <div className="modal-header">
                          <button type="button" data-dismiss="modal" aria-hidden="true" className="close">×</button>
                          <h4 id="modal-qrcode-label" className="modal-title">推广链接二维码 </h4>
                        </div>
                        <div className="modal-body">
                          <div id="pr-link-qrcode" className="text-center"></div>
                        </div>
                        <div className="modal-footer">
                            <a id="btn-download-qrcode" className="btn btn-primary" onClick={this.onClickDownloadQRCode}>下载</a>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>);        
    },
});

module.exports = {
    PrTpl: connect(mapStateToProps, mapDispatchToProps)(PrTpl)
};
