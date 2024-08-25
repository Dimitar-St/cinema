$(document).ready(function() {
    $('#signupForm').on('submit', function(event) {
        event.preventDefault(); // Prevent the default form submission

        // Get the values from the input fields
        var username = $('#username').val();
        var password = $('#password').val();

        // Perform the AJAX POST request
        $.ajax({
            url: 'localhost:8000/signup',
            type: 'POST',
            data: {
                username: username,
                password: password
            },
            success: function(response) {
                // Handle success - this can be modified according to your needs
                alert('Signup successful!');
                console.log(response);
            },
            error: function(error) {
                // Handle error - this can be modified according to your needs
                alert('Signup failed.');
                console.log(error);
            }
        });
    });
});

