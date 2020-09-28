package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/server/configurator"
	"github.com/mattermost/mattermost-plugin-apps/server/constants"

	"github.com/gorilla/mux"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-apps/server/apps"
)

const (
	SubscribePath = "/subscribe"
)

type SubscribeResponse struct {
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

type api struct {
	mm   *pluginapi.Client
	apps *apps.Service
	// subs         *apps.Subscriptions
	configurator configurator.Service
}

func Init(router *mux.Router, apps *apps.Service) {
	a := api{
		mm:           apps.Mattermost,
		configurator: apps.Config,
	}

	subrouter := router.PathPrefix(constants.APIPath).Subrouter()
	subrouter.HandleFunc(SubscribePath, a.handleSubscribe).Methods("POST", "DELETE")
}

func (a *api) handleSubscribe(w http.ResponseWriter, req *http.Request) {
	var err error

	// actingUserID := req.Header.Get("Mattermost-User-Id")
	// fmt.Printf("actingUserID = %+v\n", actingUserID)
	// if actingUserID == "" {
	// 	// err = errors.New("user not logged in")
	// 	status = http.StatusUnauthorized
	// 	return
	// }

	body, err := ioutil.ReadAll(req.Body)

	var subRequest apps.Subscription
	if err = json.Unmarshal(body, &subRequest); err != nil {
		errD, _ := json.MarshalIndent(err, "", "    ")
		fmt.Printf("err = %+v\n", string(errD))
		// return respondErr(w, http.StatusInternalServerError, err)
		return
	}

	subs := apps.NewSubscriptions(a.mm, a.configurator)

	switch req.Method {
	case http.MethodPost:
		err = subs.StoreSubscription(subRequest.Subject, subRequest, subRequest.ChannelID)
	case http.MethodDelete:
		err = subs.DeleteSubscription(subRequest.Subject, subRequest.SubscriptionID, subRequest.ChannelID)
	default:
	}
	if err != nil {
		// status = http.StatusBadRequest
		return
	}

}