function grabImage(input) {
    console.log(input.files[0]);
    if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function (e) {
            //$(input).prev().attr('src', e.target.result);
            var img = $('<img />');
            img.attr('src', e.target.result).width('100%')
            displayMessage(img[0].outerHTML, 'User');

        }
        reader.readAsDataURL(input.files[0]);
    }

    // Get form
    var form = $('#fileUploadForm')[0];
    
    // Create an FormData object 
    var data = new FormData(form);

    $.ajax({
        type: 'POST',
        enctype: 'multipart/form-data',
        url: '/simulator/receive',
        data: data,
        processData: false,
        contentType: false,
        cache: false,
        timeout: 600000,
        success: function (data) {
            try {
                var data = JSON.parse(data);
                console.log(data);
                if(data.imageURL){
                    displayMessage(`<img src = '${data.imageURL}'>`, 'JPetStore') 
                }
                displayMessage(data.response, 'JPetStore')
            } catch(e) {
                console.log('ERROR : ', e);
                displayMessage('Sorry, something went wrong.', 'JPetStore')
            }
        },
        error: function (e) {
            console.log('ERROR : ', e);
            displayMessage('Sorry, something went wrong.', 'JPetStore')
        }
    });
}

/**
 * @summary Display Chat Bubble.
 *
 * Formats the chat bubble element based on if the message is from the user or from Ana.
 *
 * @function displayMessage
 * @param {String} text - Text to be displayed in chat box.
 * @param {String} user - Denotes if the message is from Ana or the user.
 * @return null
 */
function displayMessage(text, name) {

    var chat = document.getElementById('chatBox');
    var bubble = document.createElement('div');
    bubble.className = 'message'; // Wrap the text first in a message class for common formatting

    if (name == 'JPetStore')
        bubble.innerHTML = '<div class=\'anaTitle\'>' + name + '</div><div class=\'ana\'>' + text + '</div>';
    else 
        bubble.innerHTML = '<div class=\'userTitle\'>' + name + '</div><div class=\'user\'>' + text + '</div>';

    chat.appendChild(bubble);
    chat.scrollTop = chat.scrollHeight; // Move chat down to the last message displayed
}


$( document ).ready(function() {
    displayMessage('Upload an image of a pet you are interested in!', 'JPetStore');
});