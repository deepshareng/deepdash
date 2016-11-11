var DSCOMMON = require('../lib/common.js');


var sendEmail = function() {
  var stopPoint = disabelLink(60);
  $.post('/this-is-a-clandestine-resource/resent-email', {
          username  :   emailpath
  }, function(result) {
      if(!result.success) {
          window.clearInterval(stopPoint);
          $('.resend-div').children().remove()
          addResendLink();
      }
      DSCOMMON.showModal('邮件', result.message);
  }).fail(function(err) {
      window.clearInterval(stopPoint);
      $('.resend-div').children().remove()
      addResendLink();
      DSCOMMON.showModal('邮件', "邮件发送失败");
  });
}



var addResendLink = function() {
  $('.resend-div').append($("<a href='#' class='resend-email'>").html("重新发送"))
  $('.resend-email').on('click', sendEmail);
}

addResendLink();

var disabelLink = function(second) {
  $('.resend-div').children().remove()
  $('.resend-div').html($("<a>").html(second.toString() + "秒后重发"))
  var stopPoint = setInterval(function(){
    second -= 1;
    $('.resend-div').html($("<a>").html(second.toString() + "秒后重发"))
    if (second <= 0) {
      window.clearInterval(stopPoint);
      $('.resend-div').children().remove()
      addResendLink()
    }
  }, 1000);
  return stopPoint;
}