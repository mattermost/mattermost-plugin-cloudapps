// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-apps/server/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/http/dialog"
)

func (s *service) executeInstall(params *params) (*model.CommandResponse, error) {
	println("executeInstall")
	manifestURL := ""
	appSecret := ""
	force := false
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&manifestURL, "url", "", "manifest URL")
	fs.StringVar(&appSecret, "app-secret", "", "App secret")
	fs.BoolVar(&force, "force", false, "Force re-provisioning of the app")

	err := fs.Parse(params.current)
	if err != nil {
		return normalOut(params, nil, err)
	}

	manifest, err := s.apps.Client.GetManifest(manifestURL)
	if err != nil {
		return normalOut(params, nil, err)
	}

	cc := apps.Context{}
	cc.ActingUserID = params.commandArgs.UserId
	println("before ProvisionApp")

	app, _, err := s.apps.API.ProvisionApp(
		&apps.InProvisionApp{
			ManifestURL: manifestURL,
			AppSecret:   appSecret,
			Force:       force,
		},
		&cc,
		apps.SessionToken(params.commandArgs.Session.Token),
	)
	if err != nil {
		return normalOut(params, nil, err)
	}
	println("after ProvisionApp")

	conf := s.apps.Configurator.GetConfig()

	// Finish the installation when the Dialog is submitted, see
	// <plugin>/http/dialog/install.go
	err = s.apps.Mattermost.Frontend.OpenInteractiveDialog(
		dialog.NewInstallAppDialog(manifest, appSecret, conf.PluginURL, params.commandArgs))
	if err != nil {
		return normalOut(params, nil, errors.Wrap(err, "couldn't open an interactive dialog"))
	}
	println("after OpenInteractiveDialog")

	team, err := s.apps.Mattermost.Team.Get(params.commandArgs.TeamId)
	if err != nil {
		return normalOut(params, nil, err)
	}

	return &model.CommandResponse{
		GotoLocation: params.commandArgs.SiteURL + "/" + team.Name + "/messages/@" + app.BotUsername,
		Text:         fmt.Sprintf("redirected to the DM with @%s to continue installing **%s**", app.BotUsername, manifest.DisplayName),
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	}, nil
}
