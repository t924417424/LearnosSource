let term = undefined
let socket = undefined
let upTimer = null
let downTimer = null
let resize = false
const protocol = document.location.protocol.split(':')[0];
let ws_p = "ws";
if (protocol == "https") {
    ws_p = "wss";
}
dialogDom.addEventListener('opened.mdui.dialog', function () {
    resize = false
    //console.log('opened');
    initImageConn();
    getLimit()
});

dialogDom.addEventListener('close.mdui.dialog', function () {
    resize = false
    //console.log('close');
    socket.close()
    updateIcon(runId)
    //deleteContainer()
    // term.clear()
    // term.destroy()
});

initImageConn = function(){
    socket = new WebSocket(ws_p + '://' + window.location.host + "/term/" + cid);
    term = new Terminal({cols: 180, rows: 50, screenKeys: true, cursorBlink: true, cursorStyle: "block"});
    term.open(document.getElementById('terminal'));
    fit.fit(term);
    initTerm(term, socket)
}

initTerm = function (term, socket) {
    //let msgTmp = ""
    window.onresize = function () {
        fit.fit(term);
    };
    socket.onopen = function () {
        fit.fit(term);
        let auth = {
            type: "auth",
            token: window.localStorage.getItem("token"),
        }
        socket.send(JSON.stringify(auth));  //验证权限
        //term.write("正在验证\r\n");
        //term.toggleFullscreen(true);
        term.on('data', function (data) {
            let sdata = {
                type: "cmd",
                data: data,
            }
            //console.log(isUtf8(JSON.stringify(sdata)))
            $('#term_up').addClass("mdui-text-color-green");
            upTimer = setTimeout(function () {
                $('#term_up').removeClass("mdui-text-color-green");
            },1000)
            socket.send(str2utf8(JSON.stringify(sdata)));
        });

        term.on('resize', size => {
            //console.log('resize', [size.cols, size.rows]);
            let sdata = {
                type: "resize",
                data: size.cols + "," + size.rows,
            }
            socket.send(JSON.stringify(sdata));
        });
        socket.onmessage = function (msg) {
            // console.log(term.rows)
            // console.log(rows)
            // console.log(fit.proposeGeometry(term))
            if(!resize){
                term.resize(term.cols,term.rows - 1)    //触发resize事件
                fit.fit(term);
                resize = true
            }
            // if (!is_login) {
            //     term.clear()
            //     create_sftp()
            //     is_login = true
            // }
            $('#term_down').addClass("mdui-text-color-red");
            downTimer = setTimeout(function () {
                $('#term_down').removeClass("mdui-text-color-red");
            },1000)
            let reader = new FileReader();
            reader.onload = function (event) {
                let content = reader.result;//内容就在这里
                //delete reader
                //console.log(content)
                term.write(content);
            };
            reader.readAsText(msg.data);
            //term.write(msg.data);
            //update_path(msg.data)
        };
        socket.onerror = function (e) {
            toast("连接建立出错！")
            popWindow.close()
            //is_login = false
            //layer.msg("链接出错：" + JSON.stringify(e))
            console.log(e);
        };

        socket.onclose = function (e) {
            //is_login = false
            //layer.msg("链接断开：" + JSON.stringify(e))
            //term.write("连接已断开" + "\r\n");
            toast("连接被关闭！")
            popWindow.close()
            term.clear();
            term.destroy();
        };
    };
}

function str2utf8(str) {
    let encoder = new TextEncoder('utf8');
    return encoder.encode(str);
}

function isUtf8(s) {
    let lastnames = new Array("ä", "å", "æ", "ç", "è", "é");
    let count = 0;
    for (var i = 0; i < lastnames.length; i++) {
        count += s.split(lastnames[i]).length;
    }
    if (count > s.length / 5) {
        return true;
    } else {
        return false;
    }
}