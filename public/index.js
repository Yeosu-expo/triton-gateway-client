$(function() {
    let dropdown1 = $('#dropdown-menu-model');
    let submit = $('#submit')
    let textInput = $('#text-input')
    let imageDisplay = $('#image-display')
    let loading = $('#loading')

    /* Load Model List */
    $.get('/get/model-list', function(response) {
        $.each(response, function(_, model) {
            let row = "<option value='" + model + "'>" + model + "</option>";

            dropdown1.append(row);
        });
    });

    /* Image Inference 요청 */
    $(document).ready(function() {
        submit.click(function() {
            let prompt = textInput.val();
            let model = dropdown1.val();
            let data = {
                prompt: prompt
            };

            if (model == "Select a model...") {
                alert("Please select a model");

                return;
            } else if (!prompt.trim()) {
                alert("Please enter a prompt.");

                return;
            }

            imageDisplay.empty();
            loading.text("Loading...")

            fetch("/model/" + model + "/infer", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
            .then(response => response.json())
            .then(data => {
                let base64Image = data.image;
                let imgTag = $("<img>", { src: "data:image/png;base64," + base64Image, alt: "Generated Image" });
                
                imageDisplay.empty().append(imgTag);

                loading.text("Finish");
            })
            .catch(error => {
                console.error("ERROR" + error);

                loading.text("Error")
            });
        });
    });
});