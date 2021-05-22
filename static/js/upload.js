var selectBtn;
var selectForm;
var uploadBtn;
var progressNone;
var progressDiv;

//上传进度
function uploadProgress(evt) {
    if (evt.lengthComputable) {
        var percentComplete = Math.round(evt.loaded * 100 / evt.total);
        progressDiv.css("width",percentComplete + "%").text(percentComplete + "%");
    }
    else {
        progressNone.text("Unable to compute")
    }
}

//重置
function cleanProgress(){
    progressNone.text("Progress: 0%");
    progressDiv.removeClass("progress-bar-success").css("width","0%").text("");
}

//监听选择文件信息
function fileSelected() {
    //HTML5文件API操作
    var file = selectBtn.get(0).files[0];
    if (file) {
        var fileSize = 0;
        if (file.size > 1024 * 1024)
            fileSize = (Math.round(file.size * 100 / (1024 * 1024)) / 100).toString() + 'MB';
        else
            fileSize = (Math.round(file.size * 100 / 1024) / 100).toString() + 'KB';

        $('#filename').text('Name[0]: ' + file.name);
        $('#fileSize').text('Size[0]: ' + fileSize);
        $('#fileType').text('Type[0]: ' + file.type);
    }else{
        $('#filename').text('');
        $('#fileSize').text('');
        $('#fileType').text('');
    }

    cleanProgress();
}

function upLoadFile() {
    // FormData 对象
    var form = new FormData(selectForm.get(0));
    var xhr = new XMLHttpRequest();
    progressNone.text("");
    xhr.open("post", "/upload", true);
    xhr.upload.addEventListener("progress", uploadProgress, false);
    // xhr.upload 这是html5新增的api,储存了上传过程中的信息
    xhr.upload.onprogress = function (ev) {
        var percent = 0;
        if (ev.lengthComputable) {
            percent = 100 * ev.loaded / ev.total;
            $("#myprogress").width(percent + "%");
        }
    };
    xhr.onload = function (oEvent) {
        if (xhr.status === 200) {
            progressNone.text("Upload success");
            setTimeout(function(){
                window.location.reload(true)
            },500);
        }else if(xhr.status === 400){
            alert("上传失败");
        }
    }
    xhr.send(form);
}

$(document).ready(function () {

    selectBtn = $("#userfile");
    selectForm = $("#uploadForm");
    uploadBtn = $(".btn-primary");
    progressNone = $(".progressNone");
    progressDiv = $("#progressNumber");

    //文件选择
    selectBtn.change(function () {
        fileSelected();
    })

    //文件上传
    uploadBtn.click(function () {

        if($.trim(selectBtn.val()) === ''){
            alert("请选择文件");
        }else{
            upLoadFile();
        }

    })

})