package handlers

const listTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <title>GetLive web admin</title>

    <!-- Bootstrap -->
	<!--<link href="css/bootstrap.min.css" rel="stylesheet">-->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
	<h1>GetLive web admin</h1>
	<h2>Entries</h2>
	<table class="table">
		<thead>
			<tr>
			<th scope="col">Title</th>
			<th scope="col">Time</th>
			<th scope="col">More</th>
			</tr>
		</thead>
		<tbody>
			{{ range $e := .}}
			<tr>
				<td><a href="{{ $e.URL }}" target="_blank">{{ $e.Title }}</a></td>
				<td>{{ $e.Time }}</td>
				<td>
					<a href="/entries/{{ $e.ID }}" class="btn btn-primary btn-sm">View</a>
				</td>
			</tr>
			{{ end }}
		</tbody>
	</table>

    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
  </body>
</html>`

const entryTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <title>GetLive web admin</title>

    <!-- Bootstrap -->
	<!--<link href="css/bootstrap.min.css" rel="stylesheet">-->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
	<h1>GetLive web admin</h1>
	<h2>{{ .Title }}</h2>

	<table class="table">
		<tr>
			<th>ID</th>
			<td>{{ .ID }}</td>
		</tr>
		<tr>
			<th>Time</th>
			<td>{{ .Time }}</td>
		</tr>
		<tr>
			<th>Title</th>
			<td><a href="{{ .URL }}" target="_blank">{{ .Title }}</a></td>
		</tr>
		<tr>
			<th>Description</th>
			<td>{{ .Description }}</td>
		</tr>
		<tr>
			<th>Owner</th>
			<td>{{ if .Owner }}{{ .Owner }}{{ else }}Scraped{{ end }}</td>
		</tr>
		<tr>
			<th>Approved</th>
			<td>{{ .Approved }} <a href="/entries/{{ .ID }}" class="btn btn-success btn-sm">Approve</a></td>
		</tr>
		{{ if .Approved }}
		<tr>
			<th>Approved by</th>
			<td>{{ .ApprovedBy }}</td>
		</tr>
		{{ end }}
	</table>
	
	{{ if .SocialmediaLinks }}
	<h3>Links on this page</h3>
	<table class="table">
		<thead>
			<tr>
				<th scope="col">Link</th>
			</tr>
		</thead>
		<tbody>
			{{ range $e := .SocialmediaLinks }}
			<tr>
				<td><a href="{{ $e }}" target="_blank">{{ $e }}</a></td>
			</tr>
			{{ end }}
		</tbody>
	</table>
	{{ end }}

    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
  </body>
</html>`
