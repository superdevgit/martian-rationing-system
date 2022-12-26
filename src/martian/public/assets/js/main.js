    $( document ).ready(function() {
    
    $("#packetType").on('change', function(){
        if($(this).val() == 'F'){
            $('.foodContent').show();
            $('.waterContent').hide();
        } else if (($(this).val() == 'W')) {
            $('.foodContent').hide();
            $('.waterContent').show();
        } else{
            $('.foodContent').hide();
            $('.waterContent').hide();
        }
    });

    //function validateForm() {alert('dfdf');
    $("#submit").on('click', function(){
        if ($('#packetType').val() == ''){
            $('#err').html('Please select packet type!');
            return false;
        } else if ($('#packetType').val() == 'F') {
            if ($('#packetContent').val() == ''){
                $('#err').html('Please add packet content!');
                return false;
            }
            if ($('#calories').val() == '') {
                $('#err').html('Please add food calories!');
                return false;
            }
            if ($('#expiryDate').val() == '') {
                $('#err').html('Please add expiry date!');
                return false;
            }
        } else if ($('#packetType').val() == 'W') {
            if ($('#litre').val() == '') {
                $('#err').html('Please water quantity!');
                return false;
            }
        }
        return true;
    });
    
});
