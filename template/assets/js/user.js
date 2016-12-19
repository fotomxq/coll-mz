//搜索框计时器和响应函数
var searchInputTime = '';

function searchInputKeyDown() {
    var val = $('#user-search').val();
    if (val) {
        search = val;
        page = 1;
        getList();
    }
}

//初始化page等工具组
function initTools() {
    //上一页按钮
    $('a[href="#page-prev"]').click(function() {
        if (page > 1) {
            page -= 1;
            getList();
        }
    });
    //下一页按钮
    $('a[href="#page-next"]').click(function() {
        if (userListData.length >= max) {
            page += 1;
            getList();
        }
    });
    //表格列，排序切换
    $('#user-list th').click(function() {
        if (sort != $(this).attr('data-key')) {
            sort = $(this).attr('data-key');
        } else {
            if (desc == 'true') {
                desc = 'false';
            } else {
                desc = 'true';
            }
        }
        page = 1;
        getList();
    });
    //启动搜索框计时器事件
    $('#user-search').keydown(function() {
        clearTimeout(searchInputTime);
        searchInputTime = setTimeout(searchInputKeyDown, 1000);
    });
}

//全局权限数据
var permissions = '';
//是否装载好数据
var permissionsReady = false;

//获取权限数据
function getPermissions() {
    //进入加载模式，禁止用户操作
    $('.main').addClass('loading');
    $('.main').addClass('segment');
    //从服务器获取数据
    $.get('/action-user?action=permissions', function(data) {
        if (!data) {
            return false;
        }
        if (data['status'] != true) {
            return false;
        }
        //保存数据
        permissions = data['data']['permissions'];
        permissionsReady = true;
        //写入所有权限列
        var html = '';
        for (key in permissions) {
            html += '<a class="ui label" data-key="' + key + '">' + permissions[key]['name'] + '</a>';
        }
        $('[name="permissions"]').each(function() {
            $(this).html(html);
        });
        //权限列按钮事件
        $('[name="permissions"] a').click(function() {
            $(this).parent().find('a').removeClass('blue');
            $(this).parent().attr('data-selected', $(this).attr('data-key'));
            $(this).addClass('blue');
        });
        //初始化调用列表数据
        getList();
        //初始化工具组数据
        initTools();
    }, 'json');
}

//列表全局设定
//搜索
var search = '';
//页数、页码、排序、是否倒序
var page = 1;
var max = 50;
var sort = 0;
var desc = 'true';
//用户列表数据
var userListData = '';

//获取用户列表
function getList() {
    if (permissionsReady != true) {
        return false;
    }
    //进入加载模式，禁止用户操作
    $('.main').addClass('loading');
    $('.main').addClass('segment');
    //清空数据
    $('#user-list tbody').html('');
    //从服务器获取数据
    $.get("/action-user?action=list&search=" + search + "&page=" + page + "&max=" + max + "&sort=" + sort + "&desc=" + desc, function(data) {
        if (!data) {
            return false;
        }
        if (data['status'] != true) {
            return false;
        }
        //保存数据
        userListData = data['data']['list'];
        //强行修改列表数据
        page = data['data']['page'];
        max = data['data']['max'];
        sort = data['data']['sort'];
        desc = data['data']['desc'];
        search = data['data']['search'];
        //遍历数据，生成表格数据
        var userListHtml = '';
        for (key in userListData) {
            var thisC = userListData[key];
            var date = new Date(parseInt(thisC['LastTime']) * 1000);
            var dateStr = date.getFullYear() + '-' + date.getMonth() + '-' + date.getDate() + '- ' + date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();
            var trClass = '';
            var thisStatus = '';
            if (thisC['IsDisabled'] == true) {
                trClass = 'negative';
                thisStatus = '<a class="ui label red">禁用</a> ';
            }
            userListHtml += '<tr class="' + trClass + '" data-key="' + key + '"><td>' + thisStatus + thisC['NiceName'] + '</td><td>' + thisC['UserName'] + '</td><td>' + thisC['LastIP'] + '</td><td>' + dateStr + '</td><td><div class="ui buttons mini"><a class="ui blue button" href="#list-edit"><i class="edit icon"></i>编辑</a><a class="ui yellow button" href="#list-delete"><i class="remove icon"></i>删除</a></div></td></tr>>';
        }
        //写入表格数据
        $('#user-list tbody').html(userListHtml);
        //动态样式
        $('tbody .buttons').hide();
        $('tbody tr').hover(function() {
            $(this).find('.buttons').show();
            $(this).toggleClass('warning');
        }, function() {
            $('tbody .buttons').hide();
            $(this).toggleClass('warning');
        });
        //修改页码显示数据
        $('#page-info').attr('data-text', page);
        //修改用户框架
        $('a[href="#list-edit"]').click(function() {
            var thisC = userListData[$(this).parent().parent().parent().attr('data-key')];
            $('#edit-modal').attr('data-id', thisC['ID']);
            $('#edit-user-input-nicename').val(thisC['NiceName']);
            $('#edit-user-input-username').val(thisC['UserName']);
            $('#edit-user-input-password').val('');
            $('#edit-user-input-permissions a').removeClass('blue');
            for (key in thisC['Permissions']) {
                $('#edit-user-input-permissions a[data-key="' + thisC['Permissions'][key] + '"]').addClass('blue');
            }
            if(thisC['IsDisabled']){
                $('#edit-user-input-disabled').checkbox('check');
            }else{
                $('#edit-user-input-disabled').checkbox('uncheck');
            }
            $('#edit-modal').modal('show');
        });
        $('a[href="#list-delete"]').click(function() {
            var thisC = userListData[$(this).parent().parent().parent().attr('data-key')];
            $('#delete-modal').attr('data-id', thisC['ID']);
            $('#delete-modal').modal('show');
        });
        //解除加载状态提示
        $('.main').removeClass('loading');
        $('.main').removeClass('segment');
    }, "json");
}

//添加新用户
function addUser() {
    //检查内容合法性
    var postNicename = $('#add-user-input-nicename').val();
    if (postNicename == "") {
        return false;
    }
    var postUsername = $('#add-user-input-username').val();
    if (postUsername == "") {
        return false;
    }
    var postPassword = $('#add-user-input-password').val();
    if (postPassword == "") {
        return false;
    }
    var postPasswordSha1 = hex_sha1(postPassword);
    var postPermissionArr = new Array();
    $('#add-user-input-permissions .blue').each(function() {
        postPermissionArr.push($(this).attr('data-key'));
    });
    if (postPermissionArr.length < 1) {
        return false;
    }
    var postPermissions = postPermissionArr.join('|');
    $('#add-modal').addClass('loading');
    $('#add-modal').addClass('segment');
    //提交数据
    $.post('/action-user?action=create', {
        'nicename': postNicename,
        'username': postUsername,
        'password': postPasswordSha1,
        'permissions': postPermissions
    }, function(data) {
        $('#add-modal').removeClass('loading');
        $('#add-modal').removeClass('segment');
        $('#add-modal').modal('hide');
        //反馈
        getList();
        sendMessageBool(data, "创建成功", "成功创建了用户。", "创建失败", "无法创建用户。");
    }, 'json');
}

//编辑用户
function editUser() {
    //检查内容合法性
    var userID = $('#edit-modal').attr('data-id');
    if (userID == "") {
        return false;
    }
    var postNicename = $('#edit-user-input-nicename').val();
    if (postNicename == "") {
        return false;
    }
    var postUsername = $('#edit-user-input-username').val();
    if (postUsername == "") {
        return false;
    }
    var postPassword = $('#edit-user-input-password').val();
    var postPasswordSha1 = '';
    if(postPassword != ''){
        postPasswordSha1 = hex_sha1(postPassword);
    }
    var postPermissionArr = new Array();
    $('#edit-user-input-permissions .blue').each(function() {
        postPermissionArr.push($(this).attr('data-key'));
    });
    if (postPermissionArr.length < 1) {
        return false;
    }
    var postPermissions = postPermissionArr.join('|');
    $('#edit-modal').addClass('loading');
    $('#edit-modal').addClass('segment');
    //提交数据
    $.post('/action-user?action=edit', {
        'id': userID,
        'nicename': postNicename,
        'username': postUsername,
        'password': postPasswordSha1,
        'permissions': postPermissions,
        'isdisabled' : $('#edit-user-input-disabled').checkbox('is checked')
    }, function(data) {
        $('#edit-modal').removeClass('loading');
        $('#edit-modal').removeClass('segment');
        $('#edit-modal').modal('hide');
        //反馈
        getList();
        sendMessageBool(data, "修改用户", "用户信息修改成功。", "修改失败", "无法修改该用户信息。");
    }, 'json');
}

//删除用户
function deleteUser() {
    //检查内容合法性
    var userID = $('#delete-modal').attr('data-id');
    if (userID == "") {
        return false;
    }
    $('#delete-modal').addClass('loading');
    $('#delete-modal').addClass('segment');
    //提交数据
    $.get('/action-user?action=delete&id=' + userID, function(data) {
        $('#delete-modal').removeClass('loading');
        $('#delete-modal').removeClass('segment');
        $('#delete-modal').modal('hide');
        //反馈
        getList();
        sendMessageBool(data, "删除成功", "用户已经被删除了。", "删除失败", "无法删除该用户。");
    }, 'json');
}

//初始化
$(document).ready(function() {
    //获取用户权限，该函数结束会获取用户列表数据
    getPermissions();
    //添加用户按钮
    $('a[href="#add-user"]').click(function() {
        $('#add-user-input-nicename').val('');
        $('#add-user-input-username').val('');
        $('#add-user-input-password').val('');
        $('#add-user-input-permissions a').removeClass('blue');
        $('#add-modal').modal('show');
    });
    //确认修改按钮
    $('a[href="#add-ok"]').click(function() {
        addUser();
    });
    //确认修改按钮
    $('a[href="#edit-ok"]').click(function() {
        editUser();
    });
    //确认删除按钮
    $('a[href="#delete-ok"]').click(function() {
        deleteUser();
    });
    //model关闭按钮
    $('a[href="#modal-cancel"]').click(function() {
        $(this).parent().parent().modal('hide');
    });
    //启动所有checkbox框架
    $('.ui.checkbox').checkbox();
});
