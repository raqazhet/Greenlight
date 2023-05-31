document.querySelector('form').addEventListener('submit', function(event) {
    event.preventDefault();
    const formData = {
        category: document.querySelector('input[name="category"]').value,
        text: document.querySelector('input[name="text"]').value,
    };
    fetch('/submit-form',{
        method: 'POST',
        headers: {
        'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
    })
    .then(response => {
        if (response.ok) {
            console.log('Form data suubmitted succesfully');
        } else {
            console.log('Form data submission failed');
        }
    })
    .catch(error => {
        console.error('Error submitting form data', error);
    });
});