

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
// Filter handling

class IDGen {
    constructor() {
        this.val = 0;
    }
    nextID() {
        return ++this.val
    }
}

gen = new IDGen();

var filterCritera = [];

$("#run-adv-search").click(function () {
    var queryComponents = [];
    for (var index = 0; index < filterCritera.length; index++) {
        var fieldElementID = "#" + filterCritera[index]
        queryComponents.push($(fieldElementID).children()[0].firstElementChild.value + "=" + $(fieldElementID).children()[1].firstElementChild.value)
    }
    let searchURL = "/search?" +queryComponents.join("&");
    console.log(searchURL)
    $.get(searchURL, function (data) {
        console.log(data)
    }).fail(function () {
        console.log('GET '+searchURL+' failed');
    });;


});


$("#add-filter").click(function () {
    var filterCriteraID = gen.nextID()
    console.log($("#add-filter-cirteria").val())
    switch ($("#add-filter-cirteria").val()) {
        case "by-first-name":
            $("#filters").append(getFilterByFirstNameBlock(filterCriteraID));
            filterCritera.push("fil" + filterCriteraID)
            break;
        case "by-last-name":
            $("#filters").append(filterByLastName);
            break;
        case "by-phone":
            break;
        case "by-home-phone":
            break;
        case "by-email-address":
            break;
        case "by-home-address":
            break;
        case "by-pincode":
            break;
        case "by-blood-group":
            break;
        case "by-birthday":
            break;
        case "by-anniversary":
            break;
        case "by-profession":
            break;
        case "by-age":
            break;
        case "by-marital-status":
            break;
        default:
            console.log("unhandled - add filter criteria block");
            break;
    }
    console.log(filterCritera);
});

function rem(id) {
    console.log("removing " + "fil" + id);
    $("#fil" + id).remove();
    filterCritera = filterCritera.filter(e => e !== "fil" + id);

    console.log(filterCritera);
}


function getFilterByFirstNameBlock(id) {
    return `
    <div class="form-row fil" id="fil${id}">
    <div class="col-md-3">
      <select id="gender${id}" class="form-control form-control-sm">
        <option selected value="first-name-is">First Name Is</option>
        <option  value="first-name-contains">First Name Contains</option>
        <option value="first-name-starts-with">First Name Starts With</option>
        <option value="first-name-ends-with">First Name Ends With</option>
      </select>
    </div>
    <div class="col-md-3">
      <input type="text" class="form-control  form-control-sm" id="phone" placeholder="Alex">
    </div>
    <div class="col-md-3">
      <button id="remove-filter" class="btn btn-danger btn-sm" onclick="rem(${id})"> Remove Condition</button>
    </div>
  </div>`

}


let filterByLastName = `
  <div class="form-row fil" id="filter">
  <div class="col-md-3">
    <select id="gender" class="form-control form-control-sm">
      <option selected value="last-name-is">Last Name Is</option>
      <option  value="last-name-contains">Last Name Contains</option>
      <option value="last-name-starts-with">Last Name Starts With</option>
      <option value="last-name-ends-with">Last Name Ends With</option>
    </select>
  </div>
  <div class="col-md-3">
    <input type="text" class="form-control  form-control-sm" id="phone" placeholder="Alex ">
  </div>
  <div class="col-md-3">
    <button id="remove-filter" class="btn btn-danger btn-sm"> Remove Condition</button>
  </div>
</div>`