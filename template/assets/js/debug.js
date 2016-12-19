//初始化
$(document).ready(function() {
    //表格提示
    $('tbody tr').hover(function() {
        $(this).toggleClass('warning');
    }, function() {
        $(this).toggleClass('warning');
    });
});
