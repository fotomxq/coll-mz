//从服务器传送动作
function postServerActionData(action, func) {
    $.get('/action-set?action=' + action, func);
}

//开始启动所有采集程序
function runCollAll() {
    postServerActionData('coll-all', function(data) {
        if (data == 'coll-run-ok') {
            sendNewLog('采集程序执行成功，开始获取采集日志。');
        } else {
            sendNewLog('采集程序执行失败。');
        }
    });
    //无论是否完成，都自动获取日志数据
    getCollLog();
}

//计时器和内容记录
var collLogContent = '';
var collLogTime = 1;
var collLogTimeMax = 10;

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
        //如果30次以后内容还是未变，则停止获取
        if (collLogTime > collLogTimeMax) {
            return false;
        } else {
            //每秒自动获取一次数据
            setTimeout('getCollLog()', 1000);
            //存储内容
            collLogContent = data;
        }
    });
}

//清空日志
function clearCollLog(){
    postServerActionData('clear-log', function(data) {
        clearLogContent();
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
        $(this).find('i').attr('class','setting loading icon');
        runCollAll();
    });
    //继续读取日志按钮
    $('a[href="#action-coll-log"]').click(function() {
        getCollLog();
    });
    //清空日志
    $('a[href="#action-coll-log-clear"]').click(function() {
        clearCollLog();
    });
});
