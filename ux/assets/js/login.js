function isTermOfServiceChecked(){return $("#user_agree_box")[0].checked}function isRememberMeChecked(){return $("#remember-me")[0].checked}function showApply(){return isTermOfServiceChecked()?void $("#applymodal").modal("show"):void alert("请阅读服务条款和隐私政策，并勾选“我已阅读并同意”")}function checkthebox(){$("#user_agree_box").prop("checked",!0)}function popupDiv(e){$("#"+e).animate({opacity:"show"},"slow")}function hideDiv(e){$("#"+e).animate({opacity:"hide"},"slow")}function registerUser(){registeruserurl="{{.RegisterURL}}",location.href=registeruserurl}function submitToBackEnd(){if(!isTermOfServiceChecked())return void alert("请阅读服务条款和隐私政策，并勾选“我已阅读并同意”");var e=$("#username_login").val(),o=$("#password_login").val();$.ajax({type:"POST",url:"{{.AuthURL}}",data:{username:e,password:o},success:function(n){isRememberMeChecked()?(window.localStorage.dszz=window.btoa("yes"),window.localStorage.dsxx=window.btoa(e),window.localStorage.dsyy=window.btoa(o)):(window.localStorage.dszz=window.btoa("no"),window.localStorage.removeItem("dsxx"),window.localStorage.removeItem("dsyy")),location.href="{{ .CallbackURL }}"},error:function(e){$("#wrongInfo").html("用户名或密码错误！")}})}$(document).ready(function(){$("#password_login").keyup(function(e){13==e.keyCode&&$("#login_button").click()}),$("input").click(function(){$(this).css({"background-color":"black",opacity:"0.6"})}),"yes"===window.atob(window.localStorage.dszz||"")&&($("#remember-me")[0].checked=!0,$("#username_login").val(window.atob(window.localStorage.dsxx)),$("#password_login").val(window.atob(window.localStorage.dsyy))),$("#remember-me").on("change",function(){this.checked||window.localStorage.clear()})});