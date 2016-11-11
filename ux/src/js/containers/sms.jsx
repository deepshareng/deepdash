var DSCOMMON = require('../lib/common.js');
var DSAPI = require('../lib/api.js');

var React = require('react');
var connect = require('react-redux').connect;

var mapStateToProps = function(state={}, ownProps) {
    // createStore where the reducer used
    return {
        smsInfo: state.smsInfo, 
    };
};

var mapDispatchToProps = function(dispatch) {
    return {
        initializeSmsInfo: function() {
            DSAPI.getSmsList(sessionAppId, function(data) {
                dispatch({
                    type: 'UPDATE_SMS_INFO',
                    smsInfo: data.smses,
                });
            }, function() {
                //displayNoUrlData();
            });
        },
        onClickSort: function(e) {
            var sortType = e.target.id.split('_')[1];              
            dispatch({
                type: 'UPDATE_SMS_ORDER',
                sortType: sortType,
            });
        },
    };
};

var SmsTpl = React.createClass({
    getInitialState: function() {
        return {};                
    },
    componentDidMount: function() {
        this.props.initializeSmsInfo();
        DSGLOBAL.initComponent = this.props.initializeSmsInfo;
    },
    onInputAppName: function(e) {
        if (e.which === 13) {
            this.onClickNext();
            return false;
        } 
    },
    onClickDeleteSms: function(e) {
        var self = this;
        var smsid = e.target.dataset.id;
        var confirmDelete = confirm('确认要删除推广短信吗?');

        if (!confirmDelete) {
            return;
        }

        DSAPI.deleteSms(sessionAppId, smsid, function() {
            DSCOMMON.showModal('删除推广短信', '删除推广短信成功!');
        }, function(xhr) {
            if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                DSCOMMON.showModal('删除推广短信', '没有权限，请确认账号权限或重新登录');
            } else {
                DSCOMMON.showModal('删除推广短信', '删除推广短信失败<br/>请确认网络正常，再重新尝试');
            }
        }, function() {
            self.props.initializeSmsInfo();
        });
    },
    onClickUpdateSms: function(e) {
        var self = this;
        var smsid = e.target.dataset.id;
        var content = $('#' + smsid).val();

        DSAPI.updateSms(sessionAppId, smsid, content, function() {
            DSCOMMON.showModal('修改推广短信', '修改推广短信成功');
        }, function(xhr) {
            if (xhr.responseJSON.code === DSAPI.DS_CONST_API_CODE_AUTH_FAIL) {
                DSCOMMON.showModal('修改推广短信', '没有权限，请确认账号权限或重新登录');
            } else {
                DSCOMMON.showModal('修改推广短信', '修改推广短信失败<br/>请确认网络正常，再重新尝试');
            }
        }, function() {
            self.props.initializeSmsInfo();
        });
    },
    onClickShowAddSms: function(e) {
        var self = this;
        DSAPI.getAppInfo(sessionAppId, function(data) {
            self.state.download_title = data.download_title;
            self.state.download_msg = data.download_msg;
        });
    },
    onClickAddSms: function(e) {
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
            use_shortid: true,
            channels: channels,
            download_title: self.state.download_title,
            download_msg: self.state.download_msg,
            inapp_data: getInappData(),
        };

        var content = $('#smsContent').val();

        DSAPI.generateChannelUrl(sessionAppId, param, function(data) {
            // alert success and insert 
            DSAPI.insertSms(sessionAppId, {
                'channelname': channels[0],
                'url': data.url,
                'content': content + ' ' + data.url,
            }, function(result) {
                if (!result.error) {
                    DSCOMMON.alert('生成推广短信成功!'); 
                    self.props.initializeSmsInfo();
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
        var nodes = this.props.smsInfo.smses.map(function(sms, i) {
            var names = sms.channelname.split('_');
            var cType = names[0];
            var cName = names[1];
            var cAppd = names[2];
            return (
                <tr key={sms.id}>
                    <td>{ i+1 }</td>
                    <td>{ cType }</td>
                    <td>{ cName }</td>
                    <td>{ cAppd }</td>
                    <td><textarea id={sms.id} rows="3" defaultValue={ sms.content }></textarea></td>
                    <td>
                        <a className="link btn-delete-channel" href="javascript:;" data-id={sms.id} onClick={self.onClickDeleteSms}>删除</a>
                        <a className="link btn-delete-channel" href="javascript:;" data-id={sms.id} onClick={self.onClickUpdateSms}>保存修改</a>
                    </td>
                </tr>
            );
        }); 

        if (this.props.smsInfo.smses.length <= 0) {
            $('#default-blank-channel-info').removeClass('hide');
        } else {
            $('#default-blank-channel-info').addClass('hide');
        }

        return (
        <div>
            <div className="row">
              <div className="col-lg-12">
                <h1 className="page-header">推广短信</h1>
              </div>
            </div>
            <div className="row">
              <div className="col-lg-12">
                <ul className="nav nav-tabs">
                  <li className="active"><a href="javascript:;">推广短信</a></li>
                  <div className="pull-right">
                    <div className="btn-toolbar">
                      <div className="btn-group">
                        <button type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" className="btn btn-primary btn-outline dropdown-toggle">排序 <span className="caret"></span></button>
                        <ul className="dropdown-menu">
                          <li><a id="sort_time" className="sort" onClick={this.props.onClickSort}>创建时间</a></li>
                          <li><a id="sort_name" className="sort" onClick={this.props.onClickSort}>渠道名称</a></li>
                        </ul>
                      </div>
                      <div className="btn-group">
                          <a id="btn-show-add-channel" data-toggle="modal" data-target="#myModal" className="btn btn-primary btn-outline" onClick={this.onClickShowAddSms}>添加推广短信</a>
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
                              <label htmlFor="smsContent" className="control-label">短信内容</label>
                              <textarea rows="3" id="smsContent" className="form-control"/>
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
                          <button id="btn-add-channel" type="button" className="btn btn-primary" onClick={this.onClickAddSms}>提交</button>
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
                            <th>短信内容</th>
                            <th>操作</th>
                          </tr>
                        </thead>
                        <tbody>{ nodes }</tbody>
                      </table>
                    </div>
                    <div id="default-blank-channel-info" className="text-center default-blank hide">当前应用还没有添加推广短信 ^_^</div>
                  </div>
                </div>
              </div>
            </div>
          </div>);        
    },
});

module.exports = {
    SmsTpl: connect(mapStateToProps, mapDispatchToProps)(SmsTpl)
};
