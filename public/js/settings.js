$(document).ready(function () {
    populateProfessions()
});


function populateProfessions() {
    // populating professions options
    $.get("/util/professions", function (data) {
        console.log(data)
        for (var key in data) {
            if (p.hasOwnProperty(key)) {
                console.log(key + " -> " + p[key]);
            }
        }
    }).fail(function () {
        console.log('GET /util/professions failed');
    });;
}