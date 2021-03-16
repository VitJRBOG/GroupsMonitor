function openNewAccessTokenPage() {
    var url = "/settings/access_tokens/new";
    window.location.replace(url);
}

function addAccessToken() {
    var name = document.getElementById("name").value;
	var value = document.getElementById("value").value;
    var data = {
        "name": name,
        "value": value
    };

    var url = "/settings/access_tokens/new/create";
    var redirectURL = "/settings/access_tokens";

    sendRequest("post", url, data, redirectURL);
}

function updateAccessToken(id) {
    var name = document.getElementById("name").value;
	var value = document.getElementById("value").value;
    var data = {
        "name": name,
        "value": value
    };

    var url = "/settings/access_tokens/" + id + "/update";
    var redirectURL = "/settings/access_tokens/" + id;

    sendRequest("post", url, data, redirectURL);
}

function deleteAccessToken(id) {
    var question = "Вы действительно хотите удалить запись? Действие необратимо!";
    yes = confirm(question);

    if (yes) {
        var url = "/settings/access_tokens/" + id + "/delete";
        var redirectURL = "/settings/access_tokens";
        sendRequest("post", url, {}, redirectURL);
    }
}

function addOperator() {
    var name = document.getElementById("name").value;
	var vkId = document.getElementById("vk_id").value;
    var data = {
        "name": name,
        "vk_id": vkId
    };

    var url = "/settings/operators/new/create";
    var redirectURL = "/settings/operators";

    sendRequest("post", url, data, redirectURL);
}

function updateOperator(id) {
    var name = document.getElementById("name").value;
	var vkId = document.getElementById("vk_id").value;
    var data = {
        "name": name,
        "vk_id": vkId
    };

    var url = "/settings/operators/" + id + "/update";
    var redirectURL = "/settings/operators/" + id;

    sendRequest("post", url, data, redirectURL);
}

function deleteOperator(id) {
    var question = "Вы действительно хотите удалить запись? Действие необратимо!";
    yes = confirm(question);

    if (yes) {
        var url = "/settings/operators/" + id + "/delete";
        var redirectURL = "/settings/operators";
        sendRequest("post", url, {}, redirectURL);
    }
}

function addWard() {
    var name = document.getElementById("name").value;
    var vkId = document.getElementById("vk_id").value;
    var selectorObj = document.getElementById("get_access_token");
    var getAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("post_type");
    var postType = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("wall_post_operator");
    var wallPostOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("wall_post_send_access_token");
    var wallPostSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("wall_reply_operator");
    var wallReplyOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("wall_reply_send_access_token");
    var wallReplySendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("photo_operator");
    var photoOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("photo_send_access_token");
    var photoSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("photo_comment_operator");
    var photoCommentOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("photo_comment_send_access_token");
    var photoCommentSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("video_operator");
    var videoOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("video_send_access_token");
    var videoSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("video_comment_operator");
    var videoCommentOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("video_comment_send_access_token");
    var videoCommentSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("board_post_operator");
    var boardPostOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("board_post_send_access_token");
    var boardPostSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    var data = {
        "name": name,
        "vk_id": vkId,
        "get_access_token": getAccessToken,
        "post_type": postType,
        "wall_post_operator": wallPostOperator,
        "wall_post_send_access_token": wallPostSendAccessToken,
        "wall_reply_operator": wallReplyOperator,
        "wall_reply_send_access_token": wallReplySendAccessToken,
        "photo_operator": photoOperator,
        "photo_send_access_token": photoSendAccessToken,
        "photo_comment_operator": photoCommentOperator,
        "photo_comment_send_access_token": photoCommentSendAccessToken,
        "video_operator": videoOperator,
        "video_send_access_token": videoSendAccessToken,
        "video_comment_operator": videoCommentOperator,
        "video_comment_send_access_token": videoCommentSendAccessToken,
        "board_post_operator": boardPostOperator,
        "board_post_send_access_token": boardPostSendAccessToken,
    };

    var url = "/settings/wards/new/create";
    var redirectURL = "/settings/wards";
    
    sendRequest("post", url, data, redirectURL);
}

function updateWard(id) {
    var name = document.getElementById("name").value;
    var vkId = document.getElementById("vk_id").value;
    var selectorObj = document.getElementById("get_access_token");
    var getAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("post_type");
    var postType = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("wall_post_operator");
    var wallPostOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("wall_post_send_access_token");
    var wallPostSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("wall_reply_operator");
    var wallReplyOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("wall_reply_send_access_token");
    var wallReplySendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("photo_operator");
    var photoOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("photo_send_access_token");
    var photoSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("photo_comment_operator");
    var photoCommentOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("photo_comment_send_access_token");
    var photoCommentSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("video_operator");
    var videoOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("video_send_access_token");
    var videoSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("video_comment_operator");
    var videoCommentOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("video_comment_send_access_token");
    var videoCommentSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    selectorObj = document.getElementById("board_post_operator");
    var boardPostOperator = selectorObj.options[selectorObj.selectedIndex].value;
    selectorObj = document.getElementById("board_post_send_access_token");
    var boardPostSendAccessToken = selectorObj.options[selectorObj.selectedIndex].value;

    var data = {
        "name": name,
        "vk_id": vkId,
        "get_access_token": getAccessToken,
        "post_type": postType,
        "wall_post_operator": wallPostOperator,
        "wall_post_send_access_token": wallPostSendAccessToken,
        "wall_reply_operator": wallReplyOperator,
        "wall_reply_send_access_token": wallReplySendAccessToken,
        "photo_operator": photoOperator,
        "photo_send_access_token": photoSendAccessToken,
        "photo_comment_operator": photoCommentOperator,
        "photo_comment_send_access_token": photoCommentSendAccessToken,
        "video_operator": videoOperator,
        "video_send_access_token": videoSendAccessToken,
        "video_comment_operator": videoCommentOperator,
        "video_comment_send_access_token": videoCommentSendAccessToken,
        "board_post_operator": boardPostOperator,
        "board_post_send_access_token": boardPostSendAccessToken,
    };

    var url = "/settings/wards/" + id + "/update";
    var redirectURL = "/settings/wards/" + id;
    
    sendRequest("post", url, data, redirectURL);
}

function deleteWard(id) {
    var question = "Вы действительно хотите удалить запись? Действие необратимо!";
    yes = confirm(question);

    if (yes) {
        var url = "/settings/wards/" + id + "/delete";
        var redirectURL = "/settings/wards";
        sendRequest("post", url, {}, redirectURL);
    }
}

async function sendRequest(type = '', url = '', data = {}, redirectTo = '') {
	$.ajax({
		url: url,
		type: type,
		data: data,
        success: function(result) {
            var resData = JSON.parse(result);
            if (resData["error"].length > 0) {
                document.getElementById("error_text").textContent = resData["error"];
                document.getElementById("error_msg").style.display = "block";
            } else {
                if (redirectTo.length > 0) {
                    window.location.replace(redirectTo);
                }
            }
        }
	});
}
