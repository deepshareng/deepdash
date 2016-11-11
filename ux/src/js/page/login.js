
$(document).ready(function() {
    $("#password_login").keyup(function(event) {
        if (event.keyCode == 13) {
            $("#login_button").click();
        }
    });

    $("input").focus(function() {
        /*
        $(this).css({
            "background-color": "black",
            "opacity": "0.6"
        });
        */
    }).blur(function() {
        /*
        $(this).css({
            "background-color": "#fff",
            "opacity": "1"
        });
        */
    });


    // auto complete, remember me
    if (window.atob(window.localStorage.dszz || '') === 'yes') {
        $('#remember-me')[0].checked = true;
        $("#username_login").val(window.atob(window.localStorage.dsxx));
        $("#password_login").val(window.atob(window.localStorage.dsyy));
    }

    $('#remember-me').on('change', function() {
        if (!this.checked) {
            window.localStorage.clear();
        }
    });
});

function isTermOfServiceChecked() {
    return $("#user_agree_box")[0].checked;
}
function isRememberMeChecked() {
    return $("#remember-me")[0].checked;
}

function showApply() {
    if (!isTermOfServiceChecked()) {
        alert('请阅读服务条款和隐私政策，并勾选“我已阅读并同意”');
        return;
    }

    $('#applymodal').modal('show');
}

function checkthebox() {
    $("#user_agree_box").prop("checked", true);
}

function popupDiv(div_id) {
    $("#" + div_id).animate({
        opacity: "show"
    }, "slow");
}

function hideDiv(div_id) {
    $("#" + div_id).animate({
        opacity: "hide"
    }, "slow");
}

function registerUser() {
    var registeruserurl = DS_REGISTER_URL;
    location.href = registeruserurl;

}

function submitToBackEnd() {
    if (!isTermOfServiceChecked()) {
        alert('请阅读服务条款和隐私政策，并勾选“我已阅读并同意”');
        return;
    }

    var username = $("#username_login").val();
    var password = $("#password_login").val();


    $.ajax({
        type: "POST",
        url: DS_AUTH_URL,
        data: {
            "username": username,
            "password": password,
        },
        success: function(result) {
            if (isRememberMeChecked()) {
                window.localStorage.dszz = window.btoa("yes");
                window.localStorage.dsxx = window.btoa(username);
                window.localStorage.dsyy = window.btoa(password);
            } else {
                window.localStorage.dszz = window.btoa("no");
                window.localStorage.removeItem("dsxx");
                window.localStorage.removeItem("dsyy");
            }
            if (!result.success && result.data === "unverified_user") {
                location.href = DS_RESEND_EMAIL_URL;
                return;
            }
            location.href = DS_CALLBACK_URL|| DS_LOGIN_SELECTOR_URL;
        },
        error: function(data) {
            $("#wrongInfo").html("用户名或密码错误！");
        }
    });
}

