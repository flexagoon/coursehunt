// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func IndexPage() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<hgroup><h1>CourseHunt 📚🔍</h1><h2>Find your knowledge</h2></hgroup><form hx-get=\"/search\" hx-push-url=\"true\" hx-target=\"closest main\"><fieldset role=\"search\"><input type=\"search\" placeholder=\"Search query\" name=\"q\"> <button type=\"submit\" data-loading-aria-busy>Search</button></fieldset><fieldset><label><input type=\"checkbox\" role=\"switch\" name=\"free\"> Only free courses</label> <label>Language <select name=\"language\"><option value=\"any\" selected>Any</option> <option value=\"russian\">Russian</option> <option value=\"english\">English</option></select></label> <label>Difficulty <select name=\"difficulty\"><option value=\"any\" selected>Any</option> <option value=\"beginner\">Beginner</option> <option value=\"intermediate\">Intermediate</option> <option value=\"advanced\">Advanced</option></select></label></fieldset></form><hgroup><h3>Supported sites:</h3><ul><li>Udemy</li><li>Stepik</li><li>Coursera</li><li>edX</li><li>Skillbox</li><li>ALISON</li></ul></hgroup>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
