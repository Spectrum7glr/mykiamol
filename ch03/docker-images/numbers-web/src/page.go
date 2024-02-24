package numbersweb

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Simple API Call</title>
</head>
<body>
	<h1>KIAMOL Random Number Generator</h1>
    <button id="apiCallButton">Generate Random Number!</button>
    <div id="result"></div>

    <script>
    document.getElementById("apiCallButton").addEventListener("click", function() {
        fetch('{{.ApiEndpoint}}')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // document.getElementById("result").innerText = JSON.stringify(data);
			let resultText = "Here is your random number: " + data.data;
			document.getElementById("result").innerText = resultText;
        })
        .catch(error => {
            document.getElementById("result").innerText = 'Error: ' + error.message;
        });
    });
    </script>
</body>
</html>
`
