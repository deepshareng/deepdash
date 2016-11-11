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