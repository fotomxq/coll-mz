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
        //如果菜单没创建，则创建到菜单栏
        if(!$('#menu-coll-list').html()){
            var afertHtml = '';
            for(var key in collStatusData){
                var name = collStatusData[key]["source"];
                afertHtml += '<a class="item" href="#menu-coll-list" data-key="'+name+'"><i class="selected radio icon"></i> '+name+'</a>';
            }
            $('a[href="/set"]').after('<div class="ui dropdown item" id="menu-coll-list"><div class="text"><i class="selected radio icon"></i> 采集器</div><i class="dropdown icon"></i><div class="menu">'+afertHtml+'</div></div>');
            $('#menu-coll-list').dropdown();
            $('a[href="#menu-coll-list"]').click(function(){
                setNowCollKey($(this).attr('data-key'));
            });
        }
        //如果第一次获取，初始化标签组
        if($('#coll-status').html() == ""){
            collStatusOldData = collStatusData;
            for(var key in collStatusData){
                if(collStatusData[key]['dev'] == true){
                    continue;
                }
                $('#coll-status').html($('#coll-status').html() + '<a href="#coll-status" data-key="'+key+'" class="ui grey label"><i class="selected radio icon"></i> '+collStatusData[key]['source']+'</a>');
            }
            $('a[href="#coll-status"]').click(function(){
                setNowCollKey($(this).attr('data-key'));
            });
        }
        //60秒刷新一次数据
        setTimeout('getCollStatus()', 60000);
    },'json');
}

//切换当前选择的采集器
function setNowCollKey(name){
    collNowTagKey = name;
    if(collStatusData[collNowTagKey]["status"] == true){
        sendTip("该采集器还在工作中，暂时不能浏览，请等待采集完成后再访问。");
        return false;
    }
    $('a[href="#coll-status"]').attr('class','ui grey label');
    $('a[href="#coll-status"][data-key="'+collNowTagKey+'"]').attr('class','ui blue label');
    sendBoolTip(collNowTagKey,"已经切换到" + collStatusData[collNowTagKey]["source"] + "采集器。",collStatusData[collNowTagKey]["source"]+"采集器正在运行中，请等待结束后再浏览采集数据。");
    $('#coll-content').html('');
    listParent = 0;
    $('#back-parent').hide();
    listPage = 1;
    listSearchTitle = '';
    lastData = '';
    nextPageBool = true;
    getCollView();
}

//coll view相关参数
var viewStatus = false;
var listParent = 0;
var listStar = 0;
var listSearchTitle = '';
var listPage = 1;
var listMax = 20;
var listSort = 0;
var listDesc = 'true';
var nextPageBool = true;
//避免数据重复构建
var lastData = "";
//跳转到子页面前的上级页数
var listParentPage = 1;
var listParentID = 0;

//浏览采集器内的数据
function getCollView(){
    if(nextPageBool == false){
        return false;
    }
    $.get('/action-list?coll='+collStatusData[collNowTagKey]['source']+'&parent='+listParent+'&star='+listStar+'&page='+listPage+'&title='+listSearchTitle+'&max='+listMax+'&sort='+listSort+'&desc='+listDesc, function(data){
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
            nextPageBool = false;
            $('#next-page').hide();
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
            var source = collStatusData[collNowTagKey]['source'];
            appendHtml = '<a class="column" href="#coll-node" data-source="'+source+'" data-type="'+node['file-type']+'" data-id="'+node['id']+'" data-name="'+node['name']+'">';
            if(node['name'] == "") {
                node['name'] = "无标题"
            }else{
                if(node['name'].length >= 20){
                    node['name'] = node['name'].substr(0,25) + '...';
                }
            }
            //根据类型判断插入什么内容
            switch(node['file-type']){
                case 'txt':
                    appendHtml += '<img class="ui fluid image" src="/assets/imgs/documents.png">' + node['name'];
                    break;
                case 'jpg':
                case 'gif':
                case 'jpeg':
                case 'png':
                    appendHtml += '<img class="ui fluid image" src="/action-view?coll='+source+'&id='+node['id']+'">';
                    break;
                case 'manhua-folder':
                    appendHtml += '<img class="ui fluid image" src="/assets/imgs/photos.png">' + node['name'];
                    break;
                case 'movie-folder':
                    appendHtml += '<img class="ui fluid image" src="/assets/imgs/videos.png">' + node['name'];
                    break;
                case 'folder':
                case 'txt-folder':
                case 'html-folder':
                    appendHtml += '<img class="ui fluid image" src="/assets/imgs/folder.png">' + node['name'];
                    break;
                default:
                    break;
            }
            appendHtml += '</a>';
            $('#coll-content').append(appendHtml);
        }
        //构建单击按钮
        $('a[href="#coll-node"]').click(function(){
            viewFile($(this).attr('data-source'),$(this).attr('data-type'),$(this).attr('data-id'),$(this).attr('data-name'));
        });
        //激活查询状态
        viewStatus = true;
    },'json');
}

//进入文件或文件夹
function viewFile(source,type,id,name){
    $('#show-file-title').html(name);
    imgSrc = '/action-view?coll='+source+'&id='+id;
    $('#show-file-open').attr('href',imgSrc);
    $('#show-file').attr('data-id',id);
    fileContent = '';
    switch(type){
        case 'jpg':
        case 'gif':
        case 'jpeg':
        case 'png':
            fileContent = '<img style="max-width: '+($('body').width()-100)+'px;max-height: '+($('body').height()-100)+'px;" src="'+imgSrc+'">';
            $('#show-file').dimmer('show');
            break;
        case 'txt':
        case 'mp4':
            $('#show-file').dimmer('show');
            break;
        case 'folder':
        case 'txt-folder':
        case 'manhua-folder':
        case 'movie-folder':
        case 'html-folder':
            $('#coll-content').html('');
            if(id == 0){
                listPage = listParentPage;
            }else{
                listPage = 1;
                $('#back-parent').attr('data-source',source);
                $('#back-parent').show();
            }
            listParentPage = listPage;
            listParent = id;
            listSearchTitle = '';
            lastData = '';
            nextPageBool = true;
            getCollView();
            break;
        default:
            break;
    }
    $('#show-file-content').html(fileContent);
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

//回到上一级
function backParent(){
    $('#coll-content').html('');
    listPage = listParentPage;
    listParent = listParentID;
    if(listParentID > 0){
        listParentID = 0;
        $('#back-parent').show();
    }else{
        $('#back-parent').hide();
    }
    listSearchTitle = '';
    lastData = '';
    nextPageBool = true;
    getCollView();
}

//初始化
$(document).ready(function() {
    //自适应宽度
    $('.main').css('width',$('body').width()-100+'px');
    //隐藏页码调整工具栏
    $('#page-tools').hide();
    //菜单栏强行插入页码工具显示
    var toolItems = new Array();
    toolItems[0] = {name:'<i class="chevron down icon"></i> 倒序排列',key:'desc',value:'true'};
    toolItems[1] = {name:'<i class="chevron up icon"></i> 正序排列',key:'desc',value:'false'};
    toolItems[2] = {name:'<i class="browser icon"></i> 每次获取文件数=20',key:'max',value:'20'};
    toolItems[3] = {name:'<i class="browser icon"></i> 每次获取文件数=30',key:'max',value:'30'};
    toolItems[4] = {name:'<i class="browser icon"></i> 每次获取文件数=40',key:'max',value:'40'};
    toolItemsStr = '';
    for(var key in toolItems){
        toolItemsStr += '<a class="item" href="#menu-page-tools-'+toolItems[key]['key']+'" data-value="'+toolItems[key]['value']+'">'+toolItems[key]['name']+'</a>';
    }
    $('a[href="/set"]').after('<div class="ui dropdown item"><div class="text"><i class="server icon"></i> 列表选项</div><i class="dropdown icon"></i><div class="menu">' +toolItemsStr+ '</div></div>');
    $('a[href="#menu-page-tools-desc"]').click(function(){
        listDesc = $(this).attr('data-value');
    });
    $('a[href="#menu-page-tools-max"]').click(function(){
        listMax = $(this).attr('data-value');
    });
    //在菜单栏强行插入宫格选项按钮
    var viewMode = new Array();
    viewMode[0] = {name:'宫格X6模式',value:'six'};
    viewMode[1] = {name:'宫格X4模式',value:'four'};
    viewMode[2] = {name:'宫格X3模式',value:'three'};
    viewMode[3] = {name:'宫格X2模式',value:'two'};
    viewMode[4] = {name:'宫格X1模式',value:'one'};
    viewModeHtml = '';
    for(var key in viewMode){
        viewModeHtml += '<a class="item" href="#menu-view-col" data-value="'+viewMode[key]['value']+'"><i class="eye icon"></i> '+viewMode[key]['name']+'</a>';
    }
    $('a[href="/set"]').after('<div class="ui dropdown item"><div class="text"><i class="eye icon"></i> 宫格</div><i class="dropdown icon"></i><div class="menu">' +viewModeHtml+ '</div></div>');
    //初始化所有复选框
    $('.ui.radio.checkbox').checkbox();
    //初始化所有下拉菜单
    $('.dropdown').dropdown();
    //获取status数据
    getCollStatus();
    //自动下一页
    $('#coll-content').visibility({
        once: false,
        observeChanges: true,
        onBottomVisible:function(){
            if(viewStatus == true){
                listPage += 1;
                getCollView();
            }
        }
    });
    $('#next-page').click(function(){
        if(viewStatus == true){
            listPage += 1;
            getCollView();
        }
    });
    //选项按钮宫格模式切换
    $('a[href="#menu-view-col"]').click(function(){
        modeType = $(this).attr('data-value');
        $('#coll-content').attr('class','ui '+modeType+' column grid');
    });
    //遮罩
    $('#show-file').dimmer();
    //隐藏上一级按钮
    $('#back-parent').hide();
    $('#back-parent').click(function(){
            backParent();
    });
    //关闭遮罩按钮
    $('#show-file-close').click(function(){
        $('#show-file').dimmer('hide');
    });
    //按键监听组
    $('body').keyup(function(event){
        //按下ESC，关闭遮罩
        if (event.keyCode === 27) {
            $('#show-file').dimmer('hide');
        }
        //按下BackSpace后退，回去上级
        if (event.keyCode === 8) {
            backParent();
        }
        //按下<或>方向键、空格键，切换图片
        if(event.keyCode === 37 || event.keyCode === 39 || event.keyCode === 32){
            var id = $('#show-file').attr('data-id');
            if(id < 1){
                return false;
            }
            var nextNode
            if (event.keyCode === 37) {
                nextNode = $('a[href="#coll-node"][data-id="'+id+'"]').prev();
            }
            if (event.keyCode === 39 || event.keyCode === 32) {
                nextNode = $('a[href="#coll-node"][data-id="'+id+'"]').next();
            }
            if(!nextNode.html()){
                $('#show-file').dimmer('hide');
                return false;
            }
            viewFile(nextNode.attr('data-source'),nextNode.attr('data-type'),nextNode.attr('data-id'),nextNode.attr('data-name'));
        }
    });
});
