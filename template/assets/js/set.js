//初始化启动
$(document).ready(function() {
    //初始化所有复选框
    $('.ui.radio.checkbox').checkbox();
    //初始化所有下拉菜单
    $('.ui.selection.dropdown').dropdown();
    //执行全部采集程序按钮
    $('a[href="action-coll-all"]').click(function(){
        $.get('/action-coll-run?action=coll-all',function(data){
            if(data == 'coll-run-ok'){

            }else{

            }
        })
    });
});
