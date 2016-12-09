//状态值
var collStatusData = "";
var collStatusOldData = "";
//当前单击tag
var collNowTagKey = "";
//当前采集器运行个数
var collNowRunNum = 0;

//获取运行状态
function getCollStatus(){
    //通过刷新获取数据
    $.get('/action-set?action=get-status', function(data){
        //不存在数据则返回
        if(!data){
            return false;
        }
        if(!data['login']){
            window.location.href = '/login';
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
            collStatusOldData = collStatusData;
            for(var key in collStatusData){
                $('#coll-status').html($('#coll-status').html() + '<a href="#coll-tag" data-key="'+key+'" class="ui grey label"><i class="selected radio icon"></i> '+collStatusData[key]['source']+'</a>');
            }
            $('a[href="#coll-tag"]').click(function(){
                collNowTagKey = $(this).attr('data-key');
                updateCollTip();
            });
        }
        //更新状态提示
        updateCollTip();
        //自动运行
        //如果有采集器运行，则1秒刷新
        //如果没有采集器运行，则60秒刷新
        if(collNowRunNum > 0){
            setTimeout('getCollStatus()', 1000);
        }else{
            setTimeout('getCollStatus()', 60000);
        }
    },'json');
}

//更新状态提示
function updateCollTip(){
    //采集器运行个数归零
    collNowRunNum = 0;
    //更新状态提示
    for(var key in collStatusData){
        //强制赋予当前选择为第一个key值
        if(collNowTagKey == ''){
            collNowTagKey = key;
        }
        //更新采集状态提示
        $('#coll-title').html('<p> '+collStatusData[collNowTagKey]['source'] + '采集器</p><p>地址：' + collStatusData[collNowTagKey]['url'] + '</p><p>');
        if(collStatusData[key]['status'] == true){
            $('a[href="#coll-tag"][data-key="'+collStatusData[key]['source']+'"]').attr('class','ui blue label');
            $('a[data-key="'+ key +'"]').find('i').attr('class','arrow circle down icon');
            //增加采集器运行个数计数
            collNowRunNum ++;
        }else{
            $('a[href="#coll-tag"][data-key="'+collStatusData[key]['source']+'"]').attr('class','ui grey label');
            $('a[data-key="'+key+'"]').find('i').attr('class','selected radio icon');
        }
        if(collStatusData[key]["dev"] == true){
            $('a[data-key="'+ key +'"]').find('i').attr('class','bug icon');
        }
        //判断日志数据是否超过上限，超过则尝试清空日志，并覆盖到旧日志数据中
        if(collStatusData[key]['log'].length > 10000){
            collStatusOldData[key]['log'] = collStatusData[key]['log'];
            actionCollLogClear(key);
        }else{
            if(collStatusData[key]['log'].length > 5000){
                collStatusOldData[key]['log'] = "";
            }
        }
        //显示工具按钮
        $('#coll-tools').show();
        $('#coll-tools').attr('data-key',collNowTagKey);
        //更新日志显示
        $('#log-content').html(collStatusData[collNowTagKey]['log']+collStatusOldData[collNowTagKey]['log']);
    }
    if(collNowTagKey != ""){
        if(collStatusData[collNowTagKey]['status'] == true){
            $('#coll-title').html($('#coll-title').html() + '状态：正在运行。</p>');
        }else{
            $('#coll-title').html($('#coll-title').html() + '状态：停止。</p>');
        }
        if(collStatusData[collNowTagKey]['dev']){
            $('#coll-title').html($('#coll-title').html() + '<p><div class="ui red basic label">警告，该采集器还在开发阶段，可能存在不稳定性。</div></p>');
        }
    }
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
    $("#coll-msg").css({backgroundColor:'black'});
    $("#coll-msg").animate({opacity:'0.5',color:'white'},'slow',function(){
        $("#coll-msg").css({backgroundColor:'',color:'black'});
        $("#coll-msg").animate({opacity:'1'});
    });
}

//发送b提示信号
function sendBoolTip(b,msgT,msgF) {
        if(b){
            sendTip(msgT);
        }else{
            sendTip(msgF);
        }
}

//启动某个采集器
function actionCollRun(name){
    $.get('/action-set?action=coll&name='+name, function(data){
        if(!data){
            return false;
        }
        getCollStatus();
        sendBoolTip(data['result'],'启动了'+name+'采集器。','尝试启动'+name+'采集器，但失败了。');
    },'json');
}

//关闭采集器
function actionCollClose(name){
    $.get('/action-set?action=close&name='+name, function(data){
        if(!data){
            return false;
        }
        sendBoolTip(data['result'],'强制关闭了'+name+'采集器。','尝试关闭'+name+'采集器，但失败了。');
        getCollStatus();
    },'json');
}

//清空日志
function actionCollLogClear(name){
    $.get('/action-set?action=clear-log&name='+name, function(data){
        if(!data){
            return false;
        }
        sendBoolTip(data['result'],'日志数据过于庞大，清空了'+name+'采集器日志。','尝试清空'+name+'采集器日志，但失败了。');
        getCollStatus();
    },'json');
}


//清空采集数据
function actionCollClear(name){
    $.get('/action-set?action=clear&name='+name, function(data){
        if(!data){
            return false;
        }
        sendBoolTip(data['result'],'清空了'+name+'采集器所有数据。','尝试清空'+name+'采集器数据，但失败了。');
        getCollStatus();
    },'json');
}

//初始化启动
$(document).ready(function() {
    //初始化所有复选框
    $('.ui.radio.checkbox').checkbox();
    //初始化所有下拉菜单
    $('.ui.selection.dropdown').dropdown();
    //自动启动采集器
    //runCollAll();
    //获取采集状态
    getCollStatus();
    //按钮设定
    $('#coll-tools').hide();
    $('a[href="#action-coll-run"]').click(function(){
        actionCollRun($('#coll-tools').attr('data-key'));
    });
    $('a[href="#action-coll-close"]').click(function(){
        actionCollClose($('#coll-tools').attr('data-key'));
    });
    $('a[href="#action-log-clear"]').click(function(){
        actionCollLogClear($('#coll-tools').attr('data-key'));
        collStatusOldData[$('#coll-tools').attr('data-key')]['log'] = "";
    });
    $('a[href="#action-coll-clear"]').click(function(){
        actionCollClear($('#coll-tools').attr('data-key'));
    });
});
