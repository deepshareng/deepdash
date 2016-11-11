var DSCOMMON = require('../lib/common.js');

var permissionList = {
  "channelUnlimit" : "无限制渠道数",
  "demo" : "示例权限",
};

var contains = function(list, item) {
  for (var i = 0; i < list.length; i++) {
    if (list[i] === item) {
      return true;
    }
  }
  return false;
}

var loadUserData = function(cb) {
  $.get('/this-is-a-clandestine-resource/users',{}, function(result) {
      if(!result || result.length < 1) {
        cb("加载失败")
        return;
      }
      cb(null, result) 
  }).fail(function(err) {
    cb(err, "加载失败")
  });
}

var createTableHeader = function(headList) {
  var list = ["用户名", "权限管理"];
  var thead = $("<thead></thead>");
  var tr = $("<tr></tr>");
  if (list.length < 1) {
    return thead.append(tr);
  }
  var size = "col-lg-" + Math.floor(12 / list.length);
  for(var i = 0; i < list.length; i++) {
    var th = $("<th class=" + size + "></th>").html(list[i]);
    th.appendTo(tr);
  }
  return thead.append(tr);
}

var createContent = function(userList) {
  var tbody = $("<tbody></tbody>");
  for(var i = 0; i < userList.length; i++) {
    var tr = $("<tr></tr>");
    tr.append($("<td></td>").html(userList[i].username));
    tr.append($("<td></td>").append(permissionManageBtn(userList[i].username)));
    tr.appendTo(tbody);
  }
  return tbody;
}

var permissionManageBtn = function (username) {
  return $("<button id='btn-manage_" + username + "'>").addClass("btn btn-primary btn-manage").html("权限管理");
}

var removePermissionBtn = function (id, name) {
  var button =  $("<button id=" + id + "_" + name + " type='button' class='btn btn-success removePermissionBtn'>").append($("<span aria-hidden='true'>").html(permissionList[name] + "  &times;"));
  return button;
}

var addPermissionBtn = function (id, name) {
  var button =  $("<button id=" + id + "_" + name + " type='button' class='btn btn-default addPermissionBtn'>").append($("<span aria-hidden='true'>").html(permissionList[name] + " &radic;"));
  return button;
}

var permissionBtn = function(id, name, type) {
  if (type == "add") {
    return addPermissionBtn(id, name);
  }
  if (type == "remove") {
    return removePermissionBtn(id, name);
  }
  return "";
}

var getBodyContent = function (id, list) {
  var div = $("<div class='.container-fluid'>");
  var ul = $("<ul class='list-inline'>");
  for(var item in permissionList) {
    ul.append($("<li class='margin'>").html(permissionBtn(id, item, contains(list, item) ? "remove" : "add")))
  }
  div.append(ul);
  return div
}


var createModalConent = function (userList) {
  var result = {};
  for (var i = 0; i < userList.length; i++) {
    var headContent = "更新 " + userList[i].username + " 权限";
    var bodyContent = getBodyContent(userList[i].id, userList[i].permission.split("_"));
    result[userList[i].username] = {
      head : headContent,
      body : bodyContent
    };
  }
  return result;
}

var removePermission = function () {
  $('.removePermissionBtn').on('click', function (){
    var isConfirmed = confirm('是否移除该权限？');
    var id = $(this).attr('id').split('_')[0];
    var permission = $(this).attr('id').split('_')[1];
    if (isConfirmed) {
      $.post('/this-is-a-clandestine-resource/permission/remove', {
              id   : id,
              permission : permission
      }, function(result) {
          if(!result.success) {
              alert(result.message);
          } else {
              alert(result.message);
              location.reload();
          }
      }).fail(function(err) {
          alert("移除失败")
      });
    }    
  });
}

var addPermission = function () {
  $('.addPermissionBtn').on('click', function (){
    var isConfirmed = confirm('是否添加该权限？');
    var id = $(this).attr('id').split('_')[0];
    var permission = $(this).attr('id').split('_')[1];
    if (isConfirmed) {
      $.post('/this-is-a-clandestine-resource/permission/add', {
              id   : id,
              permission : permission
      }, function(result) {
          if(!result.success) {
              alert(result.message);
          } else {
              alert(result.message);
              location.reload();
          }
      }).fail(function(err) {
          alert("添加权限失败")
      });
    }
  });
}

$(document).ready(function() {
  $('#uesr-tabel').append(createTableHeader())
  loadUserData(function(err, userList) {
    if (!!err) {
      alert("数据加载失败");
      return;
    }
    $('#uesr-tabel').append(createContent(userList));
    var content = createModalConent(userList)
    $('.btn-manage').on('click', function (argument) {
      var username = $(this).attr('id').split('_')[1];
      DSCOMMON.showModal(content[username].head,content[username].body);
      removePermission();
      addPermission();
    });
  });

})

