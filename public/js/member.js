

$(document).ready(function () {


    var t = $('#members-table').DataTable({
        data: [],
        columns: [
            { title: "#" },
            { title: "Name" },
            { title: "Birthday." },
        ]
    });


    $.get("/members", function (data) {
        if (data.length == 0) {
            reutrn
        }
        for (i = 0; i < data.length; i++) {
            console.log(data[i])
            t.row.add(["<input type=\"checkbox\" id=\"-1\" name=\"-1\" value=\"-1\">",
                "<a href=\"member-profile.html?id=" + data[i].ID + "\">" + data[i].firstname + " " + data[i].lastname + "</a>",
                data[i].birthday]).draw(false);
        }
    }).fail(function () {
        console.log('GET /members failed');
    });;
});


function createMembersTableRow(id, name, age) {
    return `<tr class="r">
    <td> <input type="checkbox" id="-1" name="-1" value="-1"> </td>
    <td>${id}</td>
    <td><a href="member-profile.html?id=${id}">${name}</a></td>
    <td>${age}</td>
  </tr>`
}

function populateMembers() {

}


$('#adv-search').click(function () {
    console.log("boom")
    if ($('#adv-search-bar').css('visibility') == 'hidden') {
        console.log("hid")
        $('#adv-search-bar').css('visibility', 'visible');
    } else {
        console.log("nope")
        $('#adv-search-bar').css('visibility', 'hidden');
    }
});


$("#checkAll").click(function () {
    var rowElems = $(".r");
    for (i = 0; i < rowElems.length; i++) {
        rowElems[i].children[0].children[0].checked = true;
    }
});

$("#checkNone").click(function () {
    var rowElems = $(".r");
    for (i = 0; i < rowElems.length; i++) {
        rowElems[i].children[0].children[0].checked = false;
    }
});

$("#deleteBtn").click(function () {
    var checkedRows = [];
    var rowElems = $(".r");
    for (i = 0; i < rowElems.length; i++) {
        if (rowElems[i].children[0].children[0].checked) {
            checkedRows.push(rowElems[i]);
        }
    }

    console.log("checkedRows: ", checkedRows);
    var delElems = [];
    for (i = 0; i < checkedRows.length; i++) {
        delElems.push(checkedRows[i].children[1].textContent);
    }
    console.log("deleting: ", delElems);

    if (confirm("Are you sure to delete " + delElems)) {
        console.log("deleting");

        for (i = 0; i < delElems.length; i++) {

            $.ajax({
                url: '/members/' + delElems[i],
                type: 'DELETE',
                success: function (result) {
                    for (i = 0; i < checkedRows.length; i++) {
                        checkedRows[i].remove();
                    }
                }
            });
        }

    } else {
        console.log("aborting ");
    }
});


