package numbersweb2

var indexPage = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Form with Token</title>
</head>
<body>
	<form action="/" method="post">
		<input type="hidden" name="token" value="{{.}}">
		<button type="submit">Call External Endpoint</button>
	</form>
</body>
</html>
`

var resultPage = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Form with Token</title>
</head>
<body>
	<h1>Response from External Endpoint</h1>
	<p>{{.Answer}}</p>
	<p>(called: {{.ApiEndpoint}})</p>
	<form action="/" method="post">
		<input type="hidden" name="token" value="{{.Token}}">
		<button type="submit">Call External Endpoint</button>
	</form>
</body>
</html>
`
var configPage = `
<!DOCTYPE html>
<html lang="en">
<head>

	<meta charset="UTF-8">
	<title>Configuration</title>
</head>
<body>
	<h1>Configuration</h1>
    <ul>
        {{range .}}
            <li>{{.}}</li>
        {{end}}
    </ul>
</body>
</html>
`
