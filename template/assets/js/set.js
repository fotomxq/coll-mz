//从服务器传送动作
function postServerActionData(action, func) {
    $.get('/action-coll-run?action=' + action, func);
}

//开始启动所有采集程序
function runCollAll() {
    postServerActionData('coll-all', function(data) {
        if (data == 'coll-run-ok') {
            sendNewLog('采集程序执行成功，开始获取采集日志。');
            getCollLog();
        } else {
            sendNewLog('采集程序执行失败。');
        }
    });
}

//计时器和内容记录
var collLogContent = '';
var collLogTime = 1;

//获取采集日志并添加到控制台
function getCollLog() {
    postServerActionData('get-log', function(data) {
        //将数据写入日志列html
        if (data) {
            clearLogContent();
            sendNewLog(data);
        }
        //如果日志内容未变，则计时器递增
        //如果日志内容改变，则计时器归零
        if (data == collLogContent) {
            collLogTime++;
        } else {
            collLogTime = 1;
        }
        //如果10秒以后内容还是未变，则停止获取
        if (collLogTime > 10) {
            return false;
        } else {
            //每秒自动获取一次数据
            setTimeout('getCollLog()', 1000);
            //存储内容
            collLogContent = data;
        }
    });
}

//向日志列发送新的日志
//倒叙陈列
function sendNewLog(msg) {
    $('#log-content').html('<p>' + msg + '</p>' + $('#log-content').html());
}

//清空日志内容
function clearLogContent() {
    $('#log-content').html('');
}

//初始化启动
$(document).ready(function() {
    //初始化所有复选框
    $('.ui.radio.checkbox').checkbox();
    //初始化所有下拉菜单
    $('.ui.selection.dropdown').dropdown();
    //执行全部采集程序按钮
    $('a[href="#action-coll-all"]').click(function() {
        runCollAll();
    });
});
