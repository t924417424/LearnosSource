let token = window.localStorage.getItem("token")
if (token) {
    window.location.href = "/console"
}

register = function () {
    let username = $('#r_username').val();
    let password = $('#r_password').val();
    let repassword = $('#r_repassword').val();
    let phone = $('#phone').val();
    let verify = $('#verifycode').val();
    if (username == "" || password == "" || phone == "" || verify == "") {
        toast("参数不能为空！")
        return false;
    }
    if (password != repassword) {
        toast("密码输入不一致！")
        return false;
    }
    $('#reg_btn').attr("disabled", true);
    http_send("/register", {username: username, password: password, phone: phone, verify: verify}, reg_callback)
}

reg_callback = function (result) {
    toast(result.msg);
    $('#reg_btn').removeAttr("disabled");
}

login = function () {
    let username = $('#username').val();
    let password = $('#password').val();
    if (username == "" || password == "") {
        toast("参数不能为空！")
        return false;
    }
    $('#login_btn').attr("disabled", true);
    http_send("/login", {username: username, password: password}, login_callback)

}

login_callback = function (result) {
    toast(result.msg);
    if (result.code == 200) {
        setTimeout(function () {
            window.location.href = "/console"
            return
        }, 2000)
    } else {
        $('#login_btn').removeAttr("disabled");
    }
}

send = function () {
    var phone = $('#phone').val();
    if (!(/^1(3|4|5|6|7|8|9)\d{9}$/.test(phone))) {
        toast("手机号码有误，请重填！");
        return false;
    }
    http_send("/send", {phone: phone}, sendmsg_callback)
}

sendmsg_callback = function (result) {
    if (result.code == 200) {
        button.attr('disabled', true);
        countdown() //按钮禁用并显示倒计时
    }
    toast(result.msg)
}


var countdowns = 130;
var button = $('#send');

countdown = function () {
    setTimeout(function () {
        if (countdowns <= 0) {
            button.removeAttr('disabled');
            button.text("发送验证");
            countdowns = 130;//60秒过后button上的文字初始化,计时器初始化;
            return;
        } else {
            button.text("(" + countdowns + "s)");
            countdowns--;
        }
        countdown()
    }, 1000) //每1000毫秒执行一次
}