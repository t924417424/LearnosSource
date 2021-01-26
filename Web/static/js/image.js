let isRun = false
let runId = 0
let popWindow = new mdui.Dialog('#dialog', {modal: true});
let dialogDom = document.getElementById('dialog');
let cid = "5547cde5-65a5-40c5-be6e-60948adb8a53";
let statusTimer = null
let limitTimer = null
let status = ["等待创建", "正在拉取镜像", "开始创建", "创建成功", "创建失败", "拉取镜像失败", "已删除", "触发资源限制"];
let statusEnum = {
    Loading: 0,
    PullImage: 1,
    StartCreate: 2,
    OkCreate: 3,
    ErrorCreate: 4,
    PullImageErr: 5,
    Deleted: 6,
    Overstep: 7,
}
// const (
// Loading      CStatus = 0
// PullImage    CStatus = 1
// StartCreate  CStatus = 2
// OkCreate     CStatus = 3
// ErrorCreate  CStatus = 4
// PullImageErr CStatus = 5
// Deleted      CStatus = 6
// Overstep     CStatus = 7
// )
//
// var msg = map[CStatus]string{
//     Loading:      "等待创建",
//         PullImage:    "正在拉取镜像",
//         PullImageErr: "拉取镜像失败",
//         StartCreate:  "开始创建",
//         OkCreate:     "创建成功",
//         ErrorCreate:  "创建失败",
//         Deleted:      "已删除",
//         Overstep:     "触发资源限制",
// }


updateIcon = function (id) {
    let btn = "#control_btn_" + id;
    let icon = "#control_ico_" + id;
    if (!isRun) {
        isRun = true
        runId = id
        createContainer(id)
        //popWindow.open()
        $(btn).removeClass("mdui-color-theme-accent");
        $(btn).addClass("mdui-color-red");
        $(icon).text("stop");
    } else {
        isRun = false
        runId = 0
        deleteContainer()
        $(btn).removeClass("mdui-color-red");
        $(btn).addClass("mdui-color-theme-accent");
        $(icon).text("play_arrow");
    }
}

createContainer = function (id) {
    loading_msg("Loading...")
    http_send("/createContainer", {id: id}, createCallback)
}

createCallback = function (result) {
    if (result.code == 200) {
        cid = result.data
        getStatus(cid)
        return
    }
    close_loading()
    toast(result.msg)
    updateIcon(runId)
}

getStatus = function (containerId) {
    loading_msg("申请资源创建")
    statusTimer = window.setInterval(function () {
        http_send("/getStatus", {cid: containerId}, statusCallback)
    }, 5000)
}

statusCallback = function (result) {
    if (result.code == 200) {
        loading_msg(status[result.data])
        if (result.data == statusEnum.OkCreate) {   //创建成功
            window.clearInterval(statusTimer)
            setTimeout(function () {
                close_loading()
            }, 1500)
            $('#containerStatus').text("创建成功")
            popWindow.open()    //打开终端窗口
        } else if (result.data == statusEnum.Deleted || result.data == statusEnum.ErrorCreate || result.data == statusEnum.PullImageErr || result.data == statusEnum.Overstep) {
            setTimeout(function () {
                close_loading()
            }, 3000)
        }
    } else {
        toast("获取容器信息失败" + result.msg)
        window.clearInterval(statusTimer)
        updateIcon(runId)
    }
    //window.clearInterval(statusTimer)
}

deleteContainer = function () {
    $('#NetWorkLimit').width("100%")
    limitTimer != null ? window.clearInterval(limitTimer) : limitTimer
    http_send("/deleteContainer", {cid: cid}, function (result) {
        if (result.code == 200) {
            toast("消息已投入删除队列")
        } else {
            toast(result.msg)
        }
    })
}

getLimit = function () {
    limitTimer = setInterval(function () {
        http_send("/getStatus", {cid: cid, limit: true}, function (result) {
            if (result.code == 200) {
                let limit = 100 - result.data + "%"
                $('#NetWorkLimit').width(limit)
                if (limit == 0) {
                    $('#containerStatus').text("触发资源限制")
                    window.clearInterval(limitTimer)
                }
            }
        })
    },15000)
}