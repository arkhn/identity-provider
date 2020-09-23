package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// authenticated checks if our cookie store has a user stored and returns the
// user's name, or an empty string if the user is not yet authenticated.
// TODO could be removed if sessions are not used
// func authenticated(r *http.Request) string {
// 	session, _ := store.Get(r, sessionName)
// 	if u, ok := session.Values["user"]; !ok {
// 		return ""
// 	} else if user, ok := u.(string); !ok {
// 		return ""
// 	} else {
// 		return user
// 	}
// }

// renderTemplate is a convenience helper for rendering templates.
func renderTemplate(w http.ResponseWriter, id string, d interface{}) bool {
	if t, err := template.New(id).ParseFiles("provider/templates/" + id); err != nil {
		http.Error(w, errors.Wrap(err, "Could not render template").Error(), http.StatusInternalServerError)
		return false
	} else if err := t.Execute(w, d); err != nil {
		http.Error(w, errors.Wrap(err, "Could not render template").Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

func parseChallengeFromRequest(r *http.Request, key string) (string, error) {
	challenge := r.URL.Query().Get(key)

	if challenge == "" {
		errMessage := fmt.Sprintf("Did not receive %s in the query.", key)
		return "", errors.New(errMessage)
	}
	return challenge, nil
}

func putAccept(url string, body interface{}, client *http.Client) string {
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
	}

	jsonResp := RedirectResp{}
	json.Unmarshal(b, &jsonResp)

	return jsonResp.RedirectTo
}
