package prom

import "net/http"

// ClientAuth is used to authenticate for client requests.
type ClientAuth struct {
	Username    string
	Password    string
	BearerToken string
}

// Apply Applies the authentication data to the request headers
func (auth *ClientAuth) Apply(req *http.Request) {
	if auth == nil {
		return
	}

	if auth.Username != "" {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	if auth.BearerToken != "" {
		token := "Bearer " + auth.BearerToken
		req.Header.Add("Authorization", token)
	}
}
