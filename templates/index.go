// Copyright (c) 2017 Femtowiki authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const indexSrc = `
{{ define "content" }}
<div id="section-tabs">
	<div id="section-search">
		<form method="GET" action="/search">
			<input type="text" name="q" placeholder="Search Femtowiki">
			<input class="btn btn-default" type="submit" value="Search">
		</form>
	</div>
	<div id="section-tabs-right">
		{{ if .IsEditMode }}
		<span class="active"><a href="{{ .EditURL }}">Source</a></span>
		<span><a href="{{ .URL }}">Read</a></span>
		{{ else }}
		<span><a href="{{ .EditURL }}">Source</a></span>
		<span class="active"><a href="{{ .URL }}">Read</a></span>
		{{ end }}
	</div>
	<div id="section-tabs-left">
		<span{{ if not .IsDiscussion }} class="active"{{ end }}><a href="{{ .URL }}">Main</a></span>
		<span{{ if .IsDiscussion }} class="active"{{ end }}><a href="{{ .URL }}?d=true">Discussion</a></span>
	</div>
</div>
<div id="meat">
{{ if .IsEditMode }}
	<form action="/editpage" method="POST">
		<input type="hidden" name="csrf" value="{{ .ctx.CSRFToken }}">
		<input type="hidden" name="t" value="{{ .cTitle }}">
		<input type="hidden" name="d" value="{{ if .IsDiscussion }}true{{ end }}">
		<textarea rows="50" name="content">{{ .Content }}</textarea>
		<input type="submit" class="btn btn-default" name="action" value="Update">
	</form>
{{ else }}
	{{ .Content }}
{{ end }}
</div>
{{ end }}`