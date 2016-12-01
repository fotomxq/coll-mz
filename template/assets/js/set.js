//从服务器传送动作
function postServerActionData(action, func) {
    $.get('/action-set?action=' + action, func);
}

//状态值
var collStatusData = "";
//当前单击tag
var collNowTagKey = "";

//获取运行状态
function getCollStatus(){
    /*
    webSocket = new WebSocket('/action-set?action=get-status');
    webSocket.onopen = function(data){

    };
    webSocket.onclose = function(data){

    };
    webSocket.onmessage = function(data){

    };
    webSocket.onerror = function(data){

    };
    */
    //通过刷新获取数据
    $.get('/action-set?action=get-status', function(data){
        //不存在数据则返回
        if(!data){
            return false;
        }
        //保存数据
        collStatusData = data;
        //如果第一次获取，初始化标签组
        if($('#coll-status').html() == ""){
            for(var key in data){
                $('#coll-status').html($('#coll-status').html() + '<a href="#coll-tag" data-key="'+key+'" class="ui grey label">'+data[key]['source']+'</a>');
            }
            $('a[href="#coll-tag"]').click(function(){
                collNowTagKey = $(this).attr('data-key');
            });
        }
        //更新状态提示
        for(var key in data){
            if(data[key]['status'] == true){
                $('a[href="#coll-tag"][data-key="'+data[key]['source']+'"]').attr('class','ui blue label');
            }else{
                $('a[href="#coll-tag"][data-key="'+data[key]['source']+'"]').attr('class','ui grey label');
            }
        }
        //变动日志显示内容
        if(collNowTagKey != ""){
            $('#log-content').html(collStatusData[collNowTagKey]['log']);
        }
    },'json');
    //自动运行
    setTimeout('getCollStatus()', 1000);
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
}

//向日志列发送新的日志
//倒叙陈列
function sendNewLog(msg) {
    $('#log-content').html('<p>' + msg + '</p>' + $('#log-content').html());
}

//初始化启动
$(document).ready(function() {
    //初始化所有复选框
    $('.ui.radio.checkbox').checkbox();
    //初始化所有下拉菜单
    $('.ui.selection.dropdown').dropdown();
    //自动启动采集器
    runCollAll();
    //获取采集状态
    getCollStatus();
});
