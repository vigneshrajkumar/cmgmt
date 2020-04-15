
function getFamilyMemberNode(id, name, isSelf, isHead) {
    var namePlaceholder = name;
    if (isHead) {
        namePlaceholder += " (Head)";
    }
    if (isSelf) {
        namePlaceholder += " (Self)";
        return "<p>" + namePlaceholder + "</p>"
    }
    return "<p><a href=\"member-profile.html?id=" + id + "\"> " + namePlaceholder + "</a></p>"
}

$(document).ready(function () {
    var memberID = getUrlParameter('id');
    console.log("fetching member of ID:", memberID);

    $.get("/members/" + memberID, function (data) {
        $("#first-name").text(data.firstname);
        $("#last-name").text(data.lastname);
        $("#birthday").text(data.birthday);
        $("#gender").text(data.gender);


        $.get("/util/member/" + data.ID + "/family", function (data) {
            console.log(data)

            for (var key in data) {
                if (data.hasOwnProperty(key)) {
                    console.log(key + " -> " + data[key]);
                    var isSelf = false;
                    if (memberID == key){
                        isSelf = true;
                    }
                    $("#family-members").append(getFamilyMemberNode(key, data[key][0], isSelf,data[key][1]));
                }
            }

            // for (var i = 0; i < data.length; ++i) {

            //     for (var key in data) {
            //         if (data.hasOwnProperty(key)) {
            //             console.log(key + " -> " + data.families[key]);
            //         }
            //     }

            //     //
            // }

        }).fail(function () {
            console.log('GET /family/{id} failed');
        });;


    }).fail(function () {
        console.log('GET /members/{id} failed');
    });;


});


var getUrlParameter = function getUrlParameter(sParam) {
    var sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
};