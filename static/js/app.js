
// Prompt is our JavaScript module for all alerts, notifications, and custom popup dialogs
function Prompt() {
    let toast = function (c) {
        const {
            title = "",
            icon = "success",
            position = "top-end"
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            text = "",
            title = "",
        } = c;
        Swal.fire({
            title: title,
            text: text,
            icon: "success",
        });
    }

    let error = function (c) {
        const {
            text = "",
            title = "",
        } = c;
        Swal.fire({
            title: title,
            text: text,
            icon: "error",
        });
    }

    async function custom(c) {
        const {
            msg = "",
            title = "",
            icon = "",
            showConfirmBtn = true
        } = c;

        const result = await Swal.fire({
            // const obj = await Swal.fire({
            title: title,
            icon: icon,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmBtn,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('startDate').value,
                    document.getElementById('endDate').value
                ]
            }
        })
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value[0] !== "" && result.value[1] !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                }
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}