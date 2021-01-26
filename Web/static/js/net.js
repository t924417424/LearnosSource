var $ = mdui.$;
GetQueryString = function (name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]);
    return null;
}

default_callback = function (result) {
    toast(result.msg)
}

let isSend = false;

http_send = function (url, data, callback = default_callback, method = "POST", limitCallback) {
    if (isSend) {
        if (limitCallback) {
            limitCallback()
        }
        return false;
    }
    url = "/api/v1" + url;
    isSend = true
    $.ajax({
        url: url,
        method: method,
        data: data,
        dataType: "json",
        beforeSend: function (xhr) {
            let token = window.localStorage.getItem("token")
            if (token) {
                xhr.setRequestHeader("Authorization", "Bearer " + token);
            }
        },
        success: function (result) {
            //layer.closeAll('loading');
            if (result.token) {
                window.localStorage.setItem("token", result.token)   //更新token
            }
            if (result.code == 301 || result.code == 302) {   //Token校验失败或过期
                logout()
                return
            }
            callback(result);
        },
        error: function (e) {
            //console.log("err")
            //layer.closeAll('loading');
            console.log(e.status);
            console.log(e.responseText)
            //layer.msg(e.responseText);
        },
        complete: function (e) {
            //layer.closeAll('loading');
            //console.log("请求完成")
            isSend = false;
        }
    });
}

toast = function (msg) {
    mdui.snackbar({
        message: msg,
        position: "top",
        timeout: 3000,
    });
}
let loading_index = null
let open_loading = null
let mask = `<div id="loading_mask" class="mdui-overlay mdui-overlay-show" style="z-index: 5100;"></div>`
loading_msg = function (msg) {
    if($('#loading_mask').length == 0){
        $('body').append(mask)
    }
    open_loading = mdui.snackbar({
        message: msg,
        position: "right-top",
        timeout: 10000,
        closeOnOutsideClick: false,
    });
    if(loading_index != null){
        close_loading()
    }
    loading_index = open_loading
    open_loading = null
}

close_loading = function(){
    if(open_loading == null){
        $('#loading_mask').remove()
    }
    loading_index.close()
}

logout = function () {
    window.localStorage.clear();        //清空浏览器缓存
    if (window != top) {
        top.location.href = "/login";
    }
    window.location.href = "/login";
}