function makeHomePage() {
    var url = "/settings/wards/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeWardObserverSwitchers(resData["wards"]);
        }
	});
}

function makeWardObserverSwitchers(wards) {
    var divObj = document.getElementById("content");

    for (let i = 0; i < wards.length; i++) {
        var divContentBoxObj = document.createElement("div");
        divContentBoxObj.className = "content_box";

        var divWardLabelObj = document.createElement("div");
        divWardLabelObj.className = "ward_label";
        divWardLabelObj.textContent = wards[i]["name"];
        divContentBoxObj.append(divWardLabelObj);

        var divWardObservingTglBtnBoxObj = document.createElement("ward_observing_tgl_btn");
        divWardObservingTglBtnBoxObj.className = "ward_observing_tgl_btn";
        divContentBoxObj.append(divWardObservingTglBtnBoxObj);

        var btnObj = document.createElement("button");
        btnObj.className = "btn_general";
        if (wards[i]["under_observation"] == 1) {
            btnObj.textContent = "Приостановить";
        } else {
            btnObj.textContent = "Запустить";
        }
        btnObj.addEventListener("click", function() {
            switchWardObservationMode(wards[i]["id"]);
        });
        divWardObservingTglBtnBoxObj.append(btnObj);

        divObj.append(divContentBoxObj);
    }
}

function makeObserverControlPage(id) {
    var url = "/observers/" + id + "/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            for (let i = 0; i < resData["wards"].length; i++) {
                var ward = resData["wards"][i]
                var isSelected = false
                if (resData["ward_id"] == ward["id"]) {
                    isSelected = true
                }
                makeRightMenuBtn("/observers/" + ward["id"], ward["name"], isSelected);
            }

            fillingUpObserverControlPage(resData);
        }
	});
}

function makeAccessTokensPage() {
    var url = "/settings/access_tokens/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeRightMenuBtn("/settings/access_tokens", "Ключи доступа", true);
            makeRightMenuUnderbtns("access_tokens", "", resData)
            makeRightMenuBtn("/settings/operators", "Операторы", false);
            makeRightMenuBtn("/settings/wards", "Подопечные", false);

            makeTableRows("access_tokens", resData["access_tokens"], "value");
        }
	});
}

function makeAccessTokenSettingsPage(id) {
    var url = "/settings/access_tokens/" + id + "/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeRightMenuBtn("/settings/access_tokens", "Ключи доступа", false);
            makeRightMenuUnderbtns("access_tokens", resData["access_token_id"], resData)
            makeRightMenuBtn("/settings/operators", "Операторы", false);
            makeRightMenuBtn("/settings/wards", "Подопечные", false);

            fillingUpAccessTokenSettingsFields(resData);
        }
	});
}

function makeOperatorsPage() {
    var url = "/settings/operators/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeRightMenuBtn("/settings/access_tokens", "Ключи доступа", false);
            makeRightMenuBtn("/settings/operators", "Операторы", true);
            makeRightMenuUnderbtns("operators", resData["operator_id"], resData)
            makeRightMenuBtn("/settings/wards", "Подопечные", false);

            makeTableRows("operators", resData["operators"], "vk_id");
        }
	});
}

function makeOperatorSettingsPage(id) {
    var url = "/settings/operators/" + id + "/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeRightMenuBtn("/settings/access_tokens", "Ключи доступа", false);
            makeRightMenuBtn("/settings/operators", "Операторы", false);
            makeRightMenuUnderbtns("operators", resData["operator_id"], resData)
            makeRightMenuBtn("/settings/wards", "Подопечные", false);

            fillingUpOperatorSettingsFields(resData);
        }
	});
}

function makeWardsPage() {
    var url = "/settings/wards/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeRightMenuBtn("/settings/access_tokens", "Ключи доступа", false);
            makeRightMenuBtn("/settings/operators", "Операторы", false);
            makeRightMenuBtn("/settings/wards", "Подопечные", true);
            makeRightMenuUnderbtns("wards", resData["ward_id"], resData)

            makeTableRows("wards", resData["wards"], "vk_id");
        }
	});
}

function makeWardGetPage() {
    var url = "/settings/wards/new/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            fillingUpWardNewFields(resData);
        }
	});
}

function makeWardSettingsPage(id) {
    var url = "/settings/wards/" + id + "/get";
    $.ajax({
		url: url,
		type: "post",
        success: function(result) {
            var resData = JSON.parse(result);

            makeRightMenuBtn("/settings/access_tokens", "Ключи доступа", false);
            makeRightMenuBtn("/settings/operators", "Операторы", false);
            makeRightMenuBtn("/settings/wards", "Подопечные", false);
            makeRightMenuUnderbtns("wards", resData["ward_id"], resData)

            fillingUpWardSettingsFields(resData);
        }
	});
}

function makeTableRows(itemsCategory, items, valueTitle) {
    for (let i = 0; i < items.length; i++) {
        var trObj = document.createElement("tr");
        trObj.className = "table_td";
        
        var tdNameObj = document.createElement("td");
        tdNameObj.className = "table_td";
        
        var aObj = document.createElement("a");
        aObj.className = "table_item_name";
        aObj.href = "/settings/" + itemsCategory + "/" + items[i]["id"];
        aObj.text = items[i]["name"];

        var tdValueObj = document.createElement("td");
        tdValueObj.className = "table_td";
        tdValueObj.style = "text-align: center;";
        tdValueObj.textContent = items[i][valueTitle];

        tdNameObj.append(aObj);
        trObj.append(tdNameObj);
        trObj.append(tdValueObj);

        var tableBody = document.getElementById(itemsCategory + "_table");
        tableBody.append(trObj);
    };
}

function makeRightMenuBtn(url, btnName, isSelected) {
    var navObj = document.getElementById("right_menu_btns");

    var pObj = document.createElement("p");

    var divObj = document.createElement("div");
    if (isSelected) {
        divObj.className = "right_menu_btn_box_selected";
    } else {
        divObj.className = "right_menu_btn_box";
    }

    var aObj = document.createElement("a");
    aObj.className = "right_menu_btn";
    aObj.href = url;
    aObj.text = btnName;

    divObj.append(aObj);

    pObj.append(divObj);

    navObj.append(pObj);
}

function makeRightMenuUnderbtns(itemsCategory, selectedItemID, resData) {
    var items = resData[itemsCategory]
    var navObj = document.getElementById("right_menu_btns");

    for (let i = 0; i < items.length; i++) {
        var pObj = document.createElement("p");

        var divObj = document.createElement("div");
        if (selectedItemID == items[i]["id"]) {
            divObj.className = "right_menu_btn_box_selected";
        } else {
            divObj.className = "right_menu_btn_box";
        }

        var aObj = document.createElement("a");
        aObj.className = "right_menu_under_btn";
        aObj.href = "/settings/" + itemsCategory + "/" + items[i]["id"];
        aObj.text = items[i]["name"];

        divObj.append(aObj);

        pObj.append(divObj);

        navObj.append(pObj);
    };
}

function fillingUpObserverControlPage(resData) {
    var divContentHeaderObj = document.getElementById("content_header");
    divContentHeaderObj.textContent += " " + resData["ward"]["name"];

    var btnWallPostSwitcherObj = document.getElementById("wall_post_switcher");
    var wallPostMode = 0;
    if (resData["lp_api_settings"]["wall_post_new"] == 1) {
        btnWallPostSwitcherObj.textContent = "Выкл.";
    } else {
        btnWallPostSwitcherObj.textContent = "Вкл.";
        wallPostMode = 1;
    }
    btnWallPostSwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "wall_post_new", wallPostMode);
    });

    var btnWallReplySwitcherObj = document.getElementById("wall_reply_switcher");
    var wallReplyMode = 0;
    if (resData["lp_api_settings"]["wall_reply_new"] == 1) {
        btnWallReplySwitcherObj.textContent = "Выкл.";
    } else {
        btnWallReplySwitcherObj.textContent = "Вкл.";
        wallReplyMode = 1;
    }
    btnWallReplySwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "wall_reply_new", wallReplyMode);
    });

    var btnPhotoSwitcherObj = document.getElementById("photo_switcher");
    var photoMode = 0;
    if (resData["lp_api_settings"]["photo_new"] == 1) {
        btnPhotoSwitcherObj.textContent = "Выкл.";
    } else {
        btnPhotoSwitcherObj.textContent = "Вкл.";
        photoMode = 1;
    }
    btnPhotoSwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "photo_new", photoMode);
    });

    var btnPhotoCommentSwitcherObj = document.getElementById("photo_comment_switcher");
    var photoCommentMode = 0;
    if (resData["lp_api_settings"]["photo_comment_new"] == 1) {
        btnPhotoCommentSwitcherObj.textContent = "Выкл.";
    } else {
        btnPhotoCommentSwitcherObj.textContent = "Вкл.";
        photoCommentMode = 1;
    }
    btnPhotoCommentSwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "photo_comment_new", photoCommentMode);
    });

    var btnVideoSwitcherObj = document.getElementById("video_switcher");
    var videoMode = 0;
    if (resData["lp_api_settings"]["video_new"] == 1) {
        btnVideoSwitcherObj.textContent = "Выкл.";
    } else {
        btnVideoSwitcherObj.textContent = "Вкл.";
        videoMode = 1;
    }
    btnVideoSwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "video_new", videoMode);
    });

    var btnVideoCommentSwitcherObj = document.getElementById("video_comment_switcher");
    var videoCommentMode = 0;
    if (resData["lp_api_settings"]["video_comment_new"] == 1) {
        btnVideoCommentSwitcherObj.textContent = "Выкл.";
    } else {
        btnVideoCommentSwitcherObj.textContent = "Вкл.";
        videoCommentMode = 1;
    }
    btnVideoCommentSwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "video_comment_new", videoCommentMode);
    });

    var btnBoardPostSwitcherObj = document.getElementById("board_post_switcher");
    var boardPostMode = 0;
    if (resData["lp_api_settings"]["board_post_new"] == 1) {
        btnBoardPostSwitcherObj.textContent = "Выкл.";
    } else {
        btnBoardPostSwitcherObj.textContent = "Вкл.";
        boardPostMode = 1;
    }
    btnBoardPostSwitcherObj.addEventListener("click", function() {
        switchObserverMode(resData["ward_id"], "board_post_new", boardPostMode);
    });
}

function fillingUpAccessTokenSettingsFields(resData) {
    var inputNameObj = document.getElementById("name");
    inputNameObj.value = resData["access_token"]["name"];
    inputNameObj.placeholder = resData["access_token"]["name"];

    var inputValueObj = document.getElementById("value");
    inputValueObj.placeholder = resData["access_token"]["value"];

    var btnSaveObj = document.getElementById("btn_save");
    btnSaveObj.addEventListener("click", function(){
        updateAccessToken(resData["access_token"]["id"]);
    })

    var btnDelObj = document.getElementById("btn_del");
    btnDelObj.addEventListener("click", function(){
        deleteAccessToken(resData["access_token"]["id"]);
    })
}

function fillingUpOperatorSettingsFields(resData) {
    var inputNameObj = document.getElementById("name");
    inputNameObj.value = resData["operator"]["name"];
    inputNameObj.placeholder = resData["operator"]["name"];

    var inputVkIDObj = document.getElementById("vk_id");
    inputVkIDObj.value = resData["operator"]["vk_id"];
    inputVkIDObj.placeholder = resData["operator"]["vk_id"];

    var btnSaveObj = document.getElementById("btn_save");
    btnSaveObj.addEventListener("click", function(){
        updateOperator(resData["operator"]["id"]);
    })

    var btnDelObj = document.getElementById("btn_del");
    btnDelObj.addEventListener("click", function(){
        deleteOperator(resData["operator"]["id"]);
    })
}

function fillingUpWardNewFields(resData) {

    var selectorGetAccessTokenObj = document.getElementById("get_access_token");
    var itemObj = new Option("", "");
    selectorGetAccessTokenObj.append(itemObj);

    for (let i = 0; i < resData["access_tokens"].length; i++) {
        let atName = resData["access_tokens"][i]["name"];
        itemObj = new Option(atName, atName);
        
        selectorGetAccessTokenObj.append(itemObj);
    }

    for (let i = 0; i < resData["observers"].length; i++) {
        let observerName = resData["observers"][i]["name"]
        if (observerName == "wall_post") {

            var selectorWallPostTypeObj = document.getElementById("post_type");
            itemObj = new Option("", "");
            selectorWallPostTypeObj.append(itemObj);

            for (let n = 0; n < resData["wall_post_types"].length; n++) {
                let itemObj = new Option(resData["wall_post_types_ru"][n], 
                    resData["wall_post_types"][n]);
                
                selectorWallPostTypeObj.append(itemObj);
            }
        }

        var selectorOperatorObj = document.getElementById(observerName + "_operator");
        var itemObj = new Option("", "");
        selectorOperatorObj.append(itemObj);

        for (let n = 0; n < resData["operators"].length; n++) {
            let operName = resData["operators"][n]["name"];

            itemObj = new Option(operName, operName);

            selectorOperatorObj.append(itemObj);
        }

        var selectorSendAccessToken = document.getElementById(observerName + "_send_access_token");
        itemObj = new Option("", "");
        selectorSendAccessToken.append(itemObj);

        for (let n = 0; n < resData["access_tokens"].length; n++) {
            let atName = resData["access_tokens"][n]["name"];

            itemObj = new Option(atName, atName);

            selectorSendAccessToken.append(itemObj);
        }
    }
}

function fillingUpWardSettingsFields(resData) {
    var inputNameObj = document.getElementById("name");
    inputNameObj.value = resData["ward"]["name"];
    inputNameObj.placeholder = resData["ward"]["name"];

    var inputVkIDObj = document.getElementById("vk_id");
    inputVkIDObj.value = resData["ward"]["vk_id"];
    inputVkIDObj.placeholder = resData["ward"]["vk_id"];

    var selectorGetAccessTokenObj = document.getElementById("get_access_token");
    for (let i = 0; i < resData["access_tokens"].length; i++) {
        let atName = resData["access_tokens"][i]["name"];
        let itemObj = new Option(atName, atName);

        if (resData["access_tokens"][i]["id"] == resData["ward"]["get_access_token_id"]) {
            itemObj.selected = true;
        }
        
        selectorGetAccessTokenObj.append(itemObj);
    }

    for (let i = 0; i < resData["observers"].length; i++) {
        let observerName = resData["observers"][i]["name"]
        if (observerName == "wall_post") {
            var selectorWallPostTypeObj = document.getElementById("post_type");

            for (let n = 0; n < resData["wall_post_types"].length; n++) {
                let itemObj = new Option(resData["wall_post_types_ru"][n], 
                    resData["wall_post_types"][n]);

                if (resData["wall_post_types"][n] == resData["observers"][i]["AdditionalParams"]["WallPost"]["post_type"]) {
                    itemObj.selected = true;
                }
                
                selectorWallPostTypeObj.append(itemObj);
            }
        }

        var selectorOperatorObj = document.getElementById(observerName + "_operator");

        for (let n = 0; n < resData["operators"].length; n++) {
            let operName = resData["operators"][n]["name"];

            let itemObj = new Option(operName, operName);

            if (resData["operators"][n]["id"] == resData["observers"][i]["operator_id"]) {
                itemObj.selected = true;
            }

            selectorOperatorObj.append(itemObj);
        }

        var selectorSendAccessToken = document.getElementById(observerName + "_send_access_token");

        for (let n = 0; n < resData["access_tokens"].length; n++) {
            let atName = resData["access_tokens"][n]["name"];

            let itemObj = new Option(atName, atName);

            if (resData["access_tokens"][n]["id"] == resData["observers"][i]["send_access_token_id"]) {
                itemObj.selected = true;
            }

            selectorSendAccessToken.append(itemObj);
        }
    }

    var btnSaveObj = document.getElementById("btn_save");
    btnSaveObj.addEventListener("click", function() {
        updateWard(resData["ward"]["id"]);
    });

    var btnDelObj = document.getElementById("btn_del");
    btnDelObj.addEventListener("click", function() {
        deleteWard(resData["ward"]["id"]);
    });
}
