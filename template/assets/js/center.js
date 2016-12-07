//状态值
var collStatusData = "";
//当前单击tag
var collNowTagKey = "";

//获取状态数据
function getCollStatus() {
    $.get('/action-set?action=get-status', function(data){
        //不存在数据则返回
        if(!data){
            return false;
        }
        if(!data['login']){
            //window.location.href = '/login';
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
                key = $(this).attr('data-key');
                if(collStatusData[key]["status"] == false){
                    collNowTagKey = key;
                }
                sendBoolTip(collStatusData[key]["status"],"已经切换到" + collStatusData[key]["source"] + "采集器。",collStatusData[key]["source"]+"采集器正在运行中，请等待结束后再浏览采集数据。");
            });
        }
        //归零
        collNowRunNum = 0;
        //更新状态提示
        for(var key in collStatusData){
            //强制赋予当前选择为第一个key值
            if(collNowTagKey == ''){
                collNowTagKey = key;
            }
        }
        //自动运行
        //如果有采集器运行，则1秒刷新
        //如果没有采集器运行，则60秒刷新
        if(collNowRunNum > 0){
            setTimeout('getCollStatus()', 1000);
        }else{
            setTimeout('getCollStatus()', 60000);
        }
    });
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

$(function() {
    //获取status数据
    getCollStatus();
});
