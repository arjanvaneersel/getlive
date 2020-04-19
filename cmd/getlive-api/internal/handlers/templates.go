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
			<th scope="col">Categories</th>
			<th scope="col">Approved</th>
			<th scope="col">More</th>
			</tr>
		</thead>
		<tbody>
			{{ range $e := .}}
			<tr>
				<td><a href="{{ $e.URL }}" target="_blank">{{ $e.Title }}</a></td>
				<td>{{ $e.Time }}</td>
				<td>{{ range $item := .Categories }}<span class="badge badge-primary">{{$item}}</span>&nbsp;{{ end }}</td>
				<td>{{ .Approved }}</td>
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
			<th>Categories</th>
			<td>{{ range $item := .Categories }}<span class="badge badge-primary">{{$item}}</span>&nbsp;{{ end }}</td>
		</tr>
		<tr>
			<th>Keywords</th>
			<td>{{ range $item := .Keywords }}<span class="badge badge-info">{{$item}}</span>&nbsp;{{ end }}</td>
		</tr>
		<tr>
			<th>Owner</th>
			<td>{{ if .Owner }}{{ .Owner }}{{ else }}Scraped{{ end }}</td>
		</tr>
		<tr>
			<th>Approved</th>
			<td>{{ .Approved }}</td>
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
	
	{{ if not .Approved }}
	<form action="/entries/{{.ID}}" method="POST">
		<div class="form-group">
			<label for="categories">Categories</label>
			<input name="categories" type="text" class="form-control" id="categories" aria-describedby="categoriesHelp" placeholder="Categories">
			<small id="categoriesHelp" class="form-text text-muted">Enter a space separated list of categories.</small>
		</div>
   		<button type="submit" class="btn btn-success">Approve</button>
	</form>
	{{ end }}

    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
  </body>
</html>`

const loginTemplate = `
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
	<h2>Login</h2>
	
	<form action="/" method="POST">
		<div class="form-group">
			<label for="email">Email</label>
			<input name="email" type="test" class="form-control" id="email" aria-describedby="emailHelp" placeholder="your@email.here">
			<small id="emailHelp" class="form-text text-muted">Enter your email address.</small>
		</div>
		<div class="form-group">
			<label for="email">Password</label>
			<input name="password" type="password" class="form-control" id="password" aria-describedby="passwordHelp" placeholder="Password">
			<small id="passwordHelp" class="form-text text-muted">Enter your password.</small>
		</div>
   		<button type="submit" class="btn btn-primary">Login</button>
	</form>

    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
  </body>
</html>`
