function readURL(input) {
    if (input.files && input.files[0]) {

        var reader = new FileReader();

        reader.onload = function (e) {
            $('.audio-upload-wrap').hide();

            $('.file-upload-content').show();

            $('.image-title').html(input.files[0].name);
        };

        reader.readAsDataURL(input.files[0]);

    } else {
        removeUpload();
    }
}

function removeUpload() {
    $('.file-upload-input').replaceWith($('.file-upload-input').clone());
    $('.file-upload-content').hide();
    $('.audio-upload-wrap').show();
}

$('.audio-upload-wrap').bind('dragover', function () {
    $('.audio-upload-wrap').addClass('audio-dropping');
});
$('.audio-upload-wrap').bind('dragleave', function () {
    $('.audio-upload-wrap').removeClass('audio-dropping');
});
