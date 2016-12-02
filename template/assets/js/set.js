//状态值
var collStatusData = "";
//当前单击tag
var collNowTagKey = "";

//获取运行状态
function getCollStatus(){
    //通过刷新获取数据
    $.get('/action-set?action=get-status', function(data){
        //不存在数据则返回
        if(!data){
            return false;
        }
        if(!data['result']){
            sendTip("无法获取状态信息。");
            return false;
        }
        //保存数据
        collStatusData = data['data'];
        //如果第一次获取，初始化标签组
        if($('#coll-status').html() == ""){
            for(var key in collStatusData){
                $('#coll-status').html($('#coll-status').html() + '<a href="#coll-tag" data-key="'+key+'" class="ui grey label">'+collStatusData[key]['source']+'</a>');
            }
            $('a[href="#coll-tag"]').click(function(){
                collNowTagKey = $(this).attr('data-key');
            });
        }
        //更新状态提示
        for(var key in collStatusData){
            if(collNowTagKey == ''){
                collNowTagKey = key
            }
            if(collStatusData[key]['status'] == true){
                $('a[href="#coll-tag"][data-key="'+collStatusData[key]['source']+'"]').attr('class','ui blue label');
            }else{
                $('a[href="#coll-tag"][data-key="'+collStatusData[key]['source']+'"]').attr('class','ui grey label');
            }
        }
        //变动日志显示内容
        if(collNowTagKey != ""){
            $('#log-content').html(collStatusData[collNowTagKey]['log']);
            $('#coll-tools').show();
            $('#coll-title').html(' ## 当前选择了'+collNowTagKey + '采集器 ## ');
            $('a[href="#action-coll-close"]').attr('data-key',collNowTagKey);
            $('a[href="#action-coll-clear"]').attr('data-key',collNowTagKey);
        }
    },'json');
    //自动运行
    setTimeout('getCollStatus()', 1000);
}

//开始启动所有采集程序
function runCollAll() {
    $.get('/action-set?action=coll&name=run-all',function(data){
        if(!data){
            return false;
        }
        sendBoolTip(data['result'],'采集程序执行成功，开始获取采集日志。','采集程序执行失败。');
    },'json');
}

//发送单一日志提示
function sendTip(msg) {
    $('#coll-msg').html(msg);
    $("#coll-msg").animate({backgroundColor:"red"});
}

//发送b提示信号
function sendBoolTip(b,msgT,msgF) {
        if(b){
            sendTip(msgT);
        }else{
            sendTip(msgF);
        }
}

//关闭采集器
function actionCollClose(name){
    $.get('/action-set?action=close&name='+name, function(data){
        if(!data){
            return false;
        }
        sendBoolTip(data['result'],'强制关闭了该采集器。','尝试关闭采集器，但失败了。');
    },'json');
}

//清空采集数据
function actionCollClear(name){
    $.get('/action-set?action=clear&name='+name, function(data){
        if(!data){
            return false;
        }
        sendBoolTip(data['result'],'清空了该采集器所有数据。','尝试清空该采集器数据，但失败了。');
    },'json');
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
    //按钮设定
    $('#coll-tools').hide();
    $('a[href="#action-coll-close"]').click(function(){
        actionCollClose($(this).attr('data-key'));
    });
    $('a[href="#action-coll-clear"]').click(function(){
        actionCollClear($(this).attr('data-key'));
    });
});
