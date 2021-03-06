// Copyright (c) 2017 Femtowiki authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const forgotpassSrc = `
{{ define "content" }}
<div class="form-container">
	<h3>Forgot password</h3>
	<form action="/forgotpass" method="POST">
		<div class="form-group">
			<input type="text" class="form-control" name="username" placeholder="Username">
		</div>
		<div class="form-group">
			<span>Remember your password?</span> <a href="/login">Signin</a>
		</div>
		{{ if .ctx.FlashMsg }}
		<div class="form-group">
			<span class="flash">{{ .ctx.FlashMsg }}</span>
		</div>
		{{ end }}
		<input type="submit" class="btn btn-default" value="Email password reset link">
	</form>
</div>
{{ end }}
`