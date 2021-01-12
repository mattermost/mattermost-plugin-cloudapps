// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/pkg/errors"
)

func LoadManifest(manifestURL string) (*api.Manifest, error) {
	var manifest api.Manifest
	resp, err := http.Get(manifestURL) // nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&manifest)
	if err != nil {
		return nil, err
	}
	err = validateManifest(&manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func validateManifest(manifest *api.Manifest) error {
	if manifest.AppID == "" {
		return errors.New("empty AppID")
	}
	if !manifest.Type.IsValid() {
		return errors.Errorf("invalid type: %s", manifest.Type)
	}

	if manifest.Type == api.AppTypeHTTP {
		_, err := url.Parse(manifest.HTTPRootURL)
		if err != nil {
			return errors.Wrapf(err, "invalid manifest URL %q", manifest.HTTPRootURL)
		}
	}
	return nil
}
