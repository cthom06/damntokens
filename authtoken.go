package main

import (
	"net/http"
	"encoding/json"
	"code.google.com/p/goauth2/oauth"
)

type resp struct {
	N string `json:"username"`
	T string `json:"damntoken"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "https://www.deviantart.com/oauth2/draft15/authorize?client_id=199&redirect_uri=http://cthom06.com/damntokens/postauth&response_type=code")
	w.WriteHeader(302)
}

func PostAuth(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("error") != "" {
		w.Write([]byte("Could not get authorization to grab token"))
		return
	}
	cfg := &oauth.Config{
		ClientId: "199",
		ClientSecret: "38dae62963bc64b06c47e3766fd9892f",
		Scope: "https://www.deviantart.com/api/draft15/",
		AuthURL: "https://www.deviantart.com/oauth2/draft15/authorize",
		TokenURL: "https://www.deviantart.com/oauth2/draft15/token",
		RedirectURL: "http://cthom06.com/postauth",
	}
	t := &oauth.Transport{Config: cfg}
	tok, err := t.Exchange(r.FormValue("code"))
	if err != nil {
		w.Write([]byte("Could not manage to grab a token because " + err.Error()))
		return
	}
	t.Token = tok
	rp, err := t.Client().Get("https://www.deviantart.com/api/draft15/user/damntoken")
	if err != nil {
		w.Write([]byte("Could not manage to grab a token because " + err.Error()))
		return
	}
	defer rp.Body.Close()
	var pk resp
	if e := json.NewDecoder(rp.Body).Decode(&pk); e != nil {
		w.Write([]byte("Could not read the token from " + pk.T + "because " + e.Error()))
		return
	}
	rp, err = t.Client().Get("https://www.deviantart.com/api/draft15/user/whoami")
	if err != nil {
		w.Write([]byte("Could not manage to grab a token because " + err.Error()))
		return
	}
	defer rp.Body.Close()
	json.NewDecoder(rp.Body).Decode(&pk)
	w.Write([]byte(`
<!DOCTYPE html>
<div style="width:640px;position:relative;top:60px;font-family:monospace;background-color:#BBBBBB;border-radius:8px;margin:0 auto; padding: 8px;text-align:center;">Authtoken for ` + pk.N + `<br /><span style="font-size:32px;">` + pk.T + `</span></div>
`))
}

func main() {
	http.HandleFunc("/damntokens/postauth", PostAuth)
	http.HandleFunc("/damntokens", Index)
	http.ListenAndServe(":3980", nil)
}
