let loading = `
    <div id="loading" class="mdui-progress">
        <div class="mdui-progress-indeterminate"></div>
    </div>
`

startLoading = function () {
    $('#appbar').append(loading)
}

stopLoading = function () {
    $('#loading').remove()
}

checkToken = function () {
    console.log("init")
    let token = window.localStorage.getItem("token")
    if (token) {
        let info = token.split(".")[1]
        info = window.atob(info)
        let time = Math.round(new Date().getTime() / 1000).toString();
        let infoObj = JSON.parse(info)
        //console.log(infoObj.exp - time - 20)
        if (infoObj.exp - time < 0) {
            logout()
            return
        } else {
            let timer = (infoObj.exp - time - 20) * 1000
            if (timer < 0) {
                timer = 1000
            }
            setTimeout(function () {
                //剩余20s时自动请求一次更新token
                http_send("/refreshToken", {}, function (r) {
                    //console.log("refreshToken")
                    checkToken()
                }, "GET", function () {
                    checkToken()        //如果有请求正在进行，则重试
                })
            }, timer)
        }
    }
}

createList = function (info) {
    return `<div class="mdui-col mdui-p-a-1">
            <div class="mdui-grid-tile">
                <img src="./static/img/bg.jpg" style="float: left"/>
                <div class="mdui-grid-tile img" style="width: 100%;height:100%;position: absolute">
                    <div class="mdui-valign" style="width: 100%;height:100%;">
                        <img class="mdui-center" src="` + info.Logo + `" style="width:30%"/>
                    </div>
                </div>
                <div class="mdui-grid-tile-actions" style="height: 10%">
                    <div class="mdui-grid-tile-text">
                        <div class="mdui-grid-tile-subtitle mdui-p-a-0">
                            <i class="mdui-icon material-icons mdui-col-xs8">cloud_circle</i>
                            ` + info.Name + `
                        </div>
                    </div>
                    <button id="control_btn_` + info.Id + `" class="mdui-btn mdui-btn-icon mdui-color-theme-accent mdui-ripple mdui-float-right mdui-shadow-3" onclick="runContainer(` + info.Id + `)">
                        <i id="control_ico_` + info.Id + `" class="mdui-icon material-icons">play_arrow</i>
                    </button>
                </div>
            </div>
        </div>`
}

runContainer = function (id) {
    if (!isRun || id == runId) {
        updateIcon(id)
    } else {
        toast("请先停止当前正使用的实例！")
    }
}

getImages = function () {
    startLoading()
    http_send("/getImage", {}, getImage_callback, "GET")
}

getImage_callback = function (result) {
    if (typeof result.data == "object" && result.code == 200) {
        for (let i in result.data) {
            $("#list").append(createList(result.data[i]))
        }
    } else {
        toast("暂无可用网关节点！")
    }
    stopLoading()
}

checkToken()
getImages()