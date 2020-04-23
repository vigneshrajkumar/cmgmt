$(document).ready(function () {
    fillMemberDetails();
});

function fillMemberDetails() {
    var memberID = getUrlParameter('id');
    console.log("fetching member of ID:", memberID);
    $.get("/members/" + memberID, function (data) {
        $("#first-name").text(data.member.firstName);
        $("#last-name").text(data.member.lastName);
        $("#phone").text(data.member.phone);
        $("#home").text(data.member.home);
        $("#email").text(data.member.email);
        $("#date-of-birth").text(data.member.dateOfBirth);
        $("#gender").text(data.member.gender);
        $("#profession").text(data.resolvedProfession);
        $("#blood-group").text(data.member.bloodGroup);
        $("#address").text(data.member.address);
        $("#pincode").text(data.member.pincode);
        $("#remarks").text(data.member.remarks);
        fillMemberRelatives(data.member.ID)
    }).fail(function () {
        console.log('GET /members/{id} failed');
    });;
}

function fillMemberRelatives(id) {
    $.get("/util/member/" + id + "/family", function (data) {
        for (var key in data) {
            if (data.hasOwnProperty(key)) {
                var isSelf = false;
                if (id == key) {
                    isSelf = true;
                }
                $("#family-members").append(getFamilyMemberNode(key, data[key][0], isSelf, data[key][1]));
            }
        }
    }).fail(function () {
        console.log('GET /family/{id} failed');
    });;
}


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