(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);throw new Error("Cannot find module '"+o+"'")}var f=n[o]={exports:{}};t[o][0].call(f.exports,function(e){var n=t[o][1][e];return s(n?n:e)},f,f.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
$('.submit').on('click', function(argument) {
  var username =  $('input[name="username"]').val()
  var password =  $('input[name="password"]').val()
  var appname  =  $('input[name="appname"]').val() 
  var phone    =  $('input[name="phone"]').val()

  var isConfirmed = confirm('是否添加新用户？');
  if (isConfirmed) {
    $.post('/this-is-a-clandestine-resource/user/add', {
            username   : username,
            password   : password,
            phone      : phone,
            appname    : appname
    }, function(result) {
        if(!result.success) {
            alert(result.message);
        } else {
            alert(result.message);
            location.reload();
        }
    }).fail(function(err) {
        alert("更新失败")
    });
  }
})
},{}]},{},[1])