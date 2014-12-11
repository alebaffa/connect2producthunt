package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Specify your configuration. (typically as a global variable)
var config = &oauth.Config{
	ClientId:     "78830bd3a0bf820629cfa4705f26ba58a4495be85678cfd3e28b017f34ea18c8",
	ClientSecret: "286634f157190e1dfd0c8730e2d833894da8066719397127d1ac24544e035af6",
	Scope:        "public",
	AuthURL:      "https://api.producthunt.com/v1/oauth/authorize",
	TokenURL:     "https://api.producthunt.com/v1/oauth/token",
	RedirectURL:  "http://localhost:3000/handler",
}

// A landing page redirects to the OAuth provider to get the auth code.
func landing(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL("foo"), http.StatusFound)
}

// The user will be redirected back to this handler, that takes the
// "code" query parameter and Exchanges it for an access token.
func handler(w http.ResponseWriter, r *http.Request) {
	t := &oauth.Transport{Config: config}
	t.Exchange(r.FormValue("code"))
	data := map[string]interface{}{}
	// The Transport now has a valid Token. Create an *http.Client
	// with which we can make authenticated API requests.
	c := t.Client()
	result, _ := c.Get("https://api.producthunt.com/v1/users")
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	json.Unmarshal(body, &data)

	fmt.Println(data)

	// btw, r.FormValue("state") == "foo"
}

func main() {
	http.HandleFunc("/", landing)
	http.HandleFunc("/handler", handler)
	http.ListenAndServe(":3000", nil)
}
