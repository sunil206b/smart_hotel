function bookRoom(roomId) {
    let html = `<form action="" id="check-availability-form" novalidate class="needs-validation" method="post">
                   <div class="form-row" id="reservationDateModel">
                       <div class="col">
                           <input disabled required class="form-control" type="text" name="check_in_date" id="check_in_date" placeholder="Checkin">
                       </div>
                       <div class="col">
                           <input disabled required class="form-control" type="text" name="check_out_date" id="check_out_date" placeholder="Checkout">
                       </div>
                   </div>
               </form>`
    attention.multiInputModel({
        msg: html,
        title: "Choose your dates!",
        willOpen: () => {
            const dateEl = document.getElementById('reservationDateModel');
            const rp = new DateRangePicker(dateEl, {
                showOnFocus: true,
                minDate: new Date(),
            })
        },
        didOpen: () => {
            document.getElementById('check_in_date').removeAttribute('disabled');
            document.getElementById('check_out_date').removeAttribute('disabled');
        },
        preConfirm: () => {
            return [
                document.getElementById('check_in_date').value,
                document.getElementById('check_out_date').value
            ]
        },
        callback: function(result) {
            let form = document.getElementById("check-availability-form");
            let formData = new FormData(form);
            formData.append("csrf_token", "{{.CSRFToken}}");
            formData.append("room_id", roomId);

            fetch('/search-availability-json', {
                method: "POST",
                body:formData,
            }).then(response => response.json()).then(data => {
                if (data.ok) {
                    let hrefURL = '/book-room?id=' + data.room_id + '&start=' + data.start_date + '&end=' + data.end_date;
                    attention.multiInputModel({
                        icon: 'success',
                        showConfirmButton: false,
                        msg: '<p>Room is available!</p>'
                            + '<p><a href="'+hrefURL+'" class="btn btn-primary">Book Now!</a>'
                    })
                } else {
                    attention.error({
                        msg: "Not Available"
                    })
                }
            })
        }
    });
}