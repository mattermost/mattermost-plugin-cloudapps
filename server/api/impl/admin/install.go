// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/utils"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

func (adm *Admin) InstallApp(cc *apps.Context, sessionToken apps.SessionToken, in *apps.InInstallApp) (*apps.App, md.MD, error) {
	// TODO <><> check if acting user is a sysadmin

	m, err := adm.store.Manifest().Get(cc.AppID)
	if err != nil {
		return nil, "", err
	}

	app, err := adm.store.App().Get(cc.AppID)
	if err != nil {
		if errors.Cause(err) != utils.ErrNotFound {
			return nil, "", err
		}
		app = &apps.App{}
	}

	app.Manifest = *m
	if app.Disabled {
		app.Disabled = false
	}
	app.GrantedPermissions = m.RequestedPermissions
	app.GrantedLocations = m.RequestedLocations
	if in.AppSecret != "" {
		app.Secret = in.AppSecret
	}

	conf := adm.conf.GetConfig()
	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	client.SetToken(string(sessionToken))

	if app.GrantedPermissions.Contains(apps.PermissionActAsBot) && app.BotUserID == "" {
		var bot *model.Bot
		var token *model.UserAccessToken
		bot, token, err = adm.ensureBot(m, cc.ActingUserID, string(sessionToken))
		if err != nil {
			return nil, "", err
		}

		app.BotUserID = bot.UserId
		app.BotUsername = bot.Username
		app.BotAccessToken = token.Token
	}

	if app.GrantedPermissions.Contains(apps.PermissionActAsUser) && app.OAuth2ClientID == "" {
		var oAuthApp *model.OAuthApp
		oAuthApp, err = adm.ensureOAuthApp(app, in.OAuth2TrustedApp, cc.ActingUserID, string(sessionToken))
		if err != nil {
			return nil, "", err
		}
		app.OAuth2ClientID = oAuthApp.Id
		app.OAuth2ClientSecret = oAuthApp.ClientSecret
		app.OAuth2TrustedApp = in.OAuth2TrustedApp
	}

	err = adm.store.App().Save(app)
	if err != nil {
		return nil, "", err
	}

	installRequest := &apps.CallRequest{
		Call:    *apps.DefaultInstallCall.WithOverrides(app.OnInstall),
		Context: cc,
	}

	resp := adm.proxy.Call(sessionToken, installRequest)
	if resp.Type == apps.CallResponseTypeError {
		return nil, "", errors.Wrap(resp, "install failed")
	}

	adm.mm.Frontend.PublishWebSocketEvent(api.WebSocketEventRefreshBindings, map[string]interface{}{}, &model.WebsocketBroadcast{UserId: cc.ActingUserID})
	return app, resp.Markdown, nil
}

func (adm *Admin) ensureOAuthApp(app *apps.App, noUserConsent bool, actingUserID, sessionToken string) (*model.OAuthApp, error) {
	conf := adm.conf.GetConfig()
	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	client.SetToken(sessionToken)

	if app.OAuth2ClientID != "" {
		oauthApp, response := client.GetOAuthApp(app.OAuth2ClientID)
		if response.StatusCode == http.StatusOK && response.Error == nil {
			_ = adm.mm.Post.DM(app.BotUserID, actingUserID, &model.Post{
				Message: fmt.Sprintf("Using existing OAuth2 App `%s`.", oauthApp.Id),
			})

			return oauthApp, nil
		}
	}

	oauth2CallbackURL := adm.conf.GetConfig().PluginURL + api.AppsPath + "/" + string(app.AppID) + api.PathOAuth2Complete

	// For the POC this should work, but for the final product I would opt for a RPC method to register the App
	oauthApp, response := client.CreateOAuthApp(&model.OAuthApp{
		CreatorId:    actingUserID,
		Name:         app.DisplayName,
		Description:  app.Description,
		CallbackUrls: []string{oauth2CallbackURL},
		Homepage:     app.HomepageURL,
		IsTrusted:    noUserConsent,
	})
	if response.StatusCode != http.StatusCreated {
		if response.Error != nil {
			return nil, errors.Wrap(response.Error, "failed to create OAuth2 App")
		}
		return nil, errors.Errorf("failed to create OAuth2 App: received status code %v", response.StatusCode)
	}

	_ = adm.mm.Post.DM(app.BotUserID, actingUserID, &model.Post{
		Message: fmt.Sprintf("Created OAuth2 App (`%s`).", oauthApp.Id),
	})

	return oauthApp, nil
}

func (adm *Admin) ensureBot(manifest *apps.Manifest, actingUserID, sessionToken string) (*model.Bot, *model.UserAccessToken, error) {
	conf := adm.conf.GetConfig()
	client := model.NewAPIv4Client(conf.MattermostSiteURL)
	client.SetToken(sessionToken)

	bot := &model.Bot{
		Username:    strings.ToLower(string(manifest.AppID)),
		DisplayName: manifest.DisplayName,
		Description: fmt.Sprintf("Bot account for `%s` App.", manifest.DisplayName),
	}

	var fullBot *model.Bot
	user, _ := client.GetUserByUsername(bot.Username, "")
	if user == nil {
		var response *model.Response
		fullBot, response = client.CreateBot(bot)

		if response.StatusCode != http.StatusCreated {
			if response.Error != nil {
				return nil, nil, response.Error
			}
			return nil, nil, errors.New("could not create bot")
		}
	} else {
		if !user.IsBot {
			return nil, nil, errors.New("a user already owns the bot username")
		}

		fullBot = model.BotFromUser(user)
		if fullBot.DeleteAt != 0 {
			var response *model.Response
			fullBot, response = client.EnableBot(fullBot.UserId)
			if response.StatusCode != http.StatusOK {
				if response.Error != nil {
					return nil, nil, response.Error
				}
				return nil, nil, errors.New("could not enable bot")
			}
		}
	}

	token, response := client.CreateUserAccessToken(fullBot.UserId, "Mattermost App Token")
	if response.StatusCode != http.StatusOK {
		if response.Error != nil {
			return nil, nil, response.Error
		}
		return nil, nil, fmt.Errorf("could not create token, status code = %v", response.StatusCode)
	}

	_ = adm.mm.Post.DM(fullBot.UserId, actingUserID, &model.Post{
		Message: fmt.Sprintf("Provisioned bot account @%s (`%s`).",
			fullBot.Username, fullBot.UserId),
	})

	return fullBot, token, nil
}
