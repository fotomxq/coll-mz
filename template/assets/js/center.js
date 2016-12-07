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
                $('#coll-status').html($('#coll-status').html() + '<a href="#coll-status" data-key="'+key+'" class="ui grey label"><i class="selected radio icon"></i> '+collStatusData[key]['source']+'</a>');
            }
            $('a[href="#coll-status"]').click(function(){
                collNowTagKey = $(this).attr('data-key');
                if(collStatusData[collNowTagKey]["status"] == true){
                    sendTip("该采集器还在工作中，暂时不能浏览，请等待采集完成后再访问。");
                    return false;
                }
                sendBoolTip(collNowTagKey,"已经切换到" + collStatusData[collNowTagKey]["source"] + "采集器。",collStatusData[collNowTagKey]["source"]+"采集器正在运行中，请等待结束后再浏览采集数据。");
                $('#coll-content').html('');
                parent = 0;
                page = 1;
                lastData = '';
                searchTitle = '';
                getCollView();
            });
        }
        //60秒刷新一次数据
        setTimeout('getCollStatus()', 60000);
    },'json');
}

//coll view相关参数
var viewStatus = false;
var parent = 0;
var star = 0;
var searchTitle = '';
var page = 1;
var max = 10;
var sort = 0;
var desc = 'false';
//避免数据重复构建
var lastData = "";

//浏览采集器内的数据
function getCollView(){
    $.get('/action-list?coll='+collStatusData[collNowTagKey]['source']+'&parent='+parent+'&star='+star+'&page='+page+'&title='+searchTitle+'&max='+max+'&sort='+sort+'&desc='+desc, function(data){
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
        if(!data['data']){
            return false;
        }
        //避免数据重复
        if(data['data'] == lastData){
            return false;
        }
        lastData = data['data']
        //遍历将数据添加到HTML中
        for(var key in data['data']){
            var node = data['data'][key];
            $('#coll-content').append('<img class="ui fluid image column" src="/action-view?coll='+collStatusData[collNowTagKey]['source']+'&id='+node['id']+'">');
        }
        //激活查询状态
        viewStatus = true;
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

//初始化
$(document).ready(function() {
    //初始化所有复选框
    $('.ui.radio.checkbox').checkbox();
    //初始化所有下拉菜单
    $('.ui.selection.dropdown').dropdown();
    //获取status数据
    getCollStatus();
    //自动下一页
    $('#coll-content').visibility({
        once: false,
        observeChanges: true,
        onBottomVisible:function(){
            if(viewStatus == true){
                page += 1;
                getCollView();
            }
        }
    });
    $('#next-page').click(function(){
        if(viewStatus == true){
            page += 1;
            getCollView();
        }
    });
});
