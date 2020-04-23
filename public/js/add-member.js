$(document).ready(function () {
    populateFamilyOptions()
    populateProfessionOptions()
});



function addMemberToExistingFamily() {
    $.ajax({
        type: "POST",
        contentType: "application/json",
        url: '/members',
        dataType: "json",
        data: JSON.stringify({

            firstName: $("#firstName").val(),
            lastName: $("#lastName").val(),

            phone: $("#phone").val(),
            home: $("#home").val(),

            email: $("#email").val(),

            dateOfBirth: $("#birthday").val(),
            gender: $("#gender").val(),

            address: $("#address").val(),
            pincode: $("#pincode").val(),

            bloodGroup: $("#blood-group").val(),
            professionID: Number($("#profession").val()),

            // photo: $("#profile-picture").val(),
            remarks: $("#remarks").val(),
            fID: $("#select-family").val()
        })
    }).done(function (data) {
        console.log("Data Loaded: " + data);
    }).fail(function () {
        console.log('POST /members failed');
    });
}

function initializeMemberIntoNewFamily() {
    console.log("initializng member into a family")
    $.ajax({
        type: "POST",
        contentType: "application/json",
        url: '/family/members',
        dataType: "json",
        data: JSON.stringify({

            firstName: $("#firstName").val(),
            lastName: $("#lastName").val(),

            phone: $("#phone").val(),
            home: $("#home").val(),

            email: $("#email").val(),

            dateOfBirth: $("#birthday").val(),
            gender: $("#gender").val(),

            address: $("#address").val(),
            pincode: $("#pincode").val(),

            bloodGroup: $("#blood-group").val(),
            professionID: Number($("#profession").val()),

            // photo: $("#profile-picture").val(),
            remarks: $("#remarks").val(),
        })
    });
}

function getFamilyOptionNode(id, name) {
    return "<option value=\"" + id + "\"> " + name + " & Fam</option>"
}

function getProfessionOptionNode(id, name) {
    return "<option value=\"" + id + "\"> " + name + " </option>"
}

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



function populateFamilyOptions() {
    $.get("/util/families", function (data) {
        for (var key in data.families) {
            if (data.families.hasOwnProperty(key)) {
                $("#select-family").append(getFamilyOptionNode(key, data.families[key]));
            }
        }
    }).fail(function () {
        console.log('GET /util/families failed');
    });
}


function populateProfessionOptions() {
    // populating profession options
    $.get("/util/professions", function (data) {
        for (var key in data) {
            if (data.hasOwnProperty(key)) {
                $("#profession").append(getProfessionOptionNode(key, data[key]));
            }
        }
    }).fail(function () {
        console.log('GET /util/professions failed');
    });;
}