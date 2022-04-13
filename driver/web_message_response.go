package driver

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/ory/fosite"
)

var WebMessageDefaultTemplate = template.Must(template.New("web_message").Parse(`<html>
<body>
<script type="text/javascript">
var callbackOrigin = "{{ .RedirURL }}";

var response = {
{{ range $key,$value := .Parameters }}
	{{ range $parameter:= $value}}
	"{{$key}}": "{{$parameter}}",
	{{end}}
{{ end }}
};

var authorizationResponse = {
	type: "authorization_response",
	response
};


(function(window, document) {
	const isPopup = window.opener && window.opener !== window;
	if (isPopup) {
		// we're in a popup
		window.opener.postMessage(authorizationResponse, callbackOrigin);
	} else {
		// we're in a frame
		window.parent.postMessage(authorizationResponse, callbackOrigin);
	}
}
)(this, this.document);
</script>
</body>
</html>
`))

type WebMessageResponse struct {
}

func (d *WebMessageResponse) ResponseModes() fosite.ResponseModeTypes {
	return fosite.ResponseModeTypes{"web_message"}
}
func (d *WebMessageResponse) WriteAuthorizeResponse(rw http.ResponseWriter, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) {
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
	err := WebMessageDefaultTemplate.Execute(rw, struct {
		RedirURL   string
		Parameters url.Values
	}{
		RedirURL:   ar.GetRedirectURI().String(),
		Parameters: resp.GetParameters(),
	})
	if err != nil {
		log.Print("[WARN] Failed to write web message response: ", err)
	}
}
func (d *WebMessageResponse) WriteAuthorizeError(rw http.ResponseWriter, ar fosite.AuthorizeRequester, err error) {
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
	e := WebMessageDefaultTemplate.Execute(rw, struct {
		RedirURL   string
		Parameters url.Values
	}{
		RedirURL: ar.GetRedirectURI().String(),
		Parameters: url.Values{
			"error": []string{err.Error()},
		},
	})
	if e != nil {
		log.Print("[WARN] Failed to write web message response: ", err)
	}
}
