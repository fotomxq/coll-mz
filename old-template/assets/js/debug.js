$(document).ready(function() {
    $('#debug-content').css('width',$('.container').width()+"px");
    $('#debug-content').css('height',$('.container').height()+"px");
    $('a[href="#coll-debug"]').click(function(){
        $.get($(this).attr('data-url'),function(data){
            if(!data){
                return false;
            }
            $('#debug-content').html(data['html']);
        },'json')
    });
});