function addMemberToExistingFamily() {
    $.post("/members", {
        firstname: $("#firstName").val(),
        lastname: $("#lastName").val(),
        birthday: $("#birthday").val(),
        gender: $("#gender").val(),
        fID: $("#select-family").val()
    })
        .done(function (data) {
            console.log("Data Loaded: " + data);
        });
}

function initializeMemberIntoNewFamily() {
    $.post("/family/members", {
        firstname: $("#firstName").val(),
        lastname: $("#lastName").val(),
        birthday: $("#birthday").val(),
        gender: $("#gender").val(),
    })
        .done(function (data) {
            console.log("Data Loaded: " + data);
        });
}

function getOptionNode(id, name){
    return "<option value=\""+id+"\"> " + name + " & Fam</option>"
}

$(document).ready(function () {
    // Toggling family attach
    $("#init-family").change(function () {
        if (this.checked) {
            $("#select-family-grp").css('visibility', 'hidden');
        } else {
            $("#select-family-grp").css('visibility', 'visible');
        }
    });

    // Add Member
    $("#add-member").click(function () {
        if ($("#init-family").is(':checked')) {
            console.log("checked: ")
            initializeMemberIntoNewFamily()
        } else {
            console.log("un checked: ")
            addMemberToExistingFamily()
        }
    });

    // populating select-family options
    $.get("/util/families", function (data) {
        console.log(data);

        for (var key in data.families) {
            if (data.families.hasOwnProperty(key)) {
                console.log(key + " -> " + data.families[key]);
                $("#select-family").append(getOptionNode(key, data.families[key]));
            }
        }

    }).fail(function () {
        console.log('GET /families failed');
    });;

});





