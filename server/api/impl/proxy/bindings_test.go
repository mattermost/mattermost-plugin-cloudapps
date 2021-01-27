package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-apps/server/api"
	"github.com/mattermost/mattermost-plugin-apps/server/api/mock_api"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testBinding(appID api.AppID, parent api.Location, n string) []*api.Binding {
	return []*api.Binding{
		{
			AppID:    appID,
			Location: parent,
			Bindings: []*api.Binding{
				{
					AppID:    appID,
					Location: api.Location(fmt.Sprintf("id-%s", n)),
					Hint:     fmt.Sprintf("hint-%s", n),
				},
			},
		},
	}
}

func TestMergeBindings(t *testing.T) {
	type TC struct {
		name               string
		bb1, bb2, expected []*api.Binding
	}

	for _, tc := range []TC{
		{
			name: "happy simplest",
			bb1: []*api.Binding{
				{
					Location: "1",
				},
			},
			bb2: []*api.Binding{
				{
					Location: "2",
				},
			},
			expected: []*api.Binding{
				{
					Location: "1",
				},
				{
					Location: "2",
				},
			},
		},
		{
			name:     "happy simple 1",
			bb1:      testBinding("app1", api.LocationCommand, "simple"),
			bb2:      nil,
			expected: testBinding("app1", api.LocationCommand, "simple"),
		},
		{
			name:     "happy simple 2",
			bb1:      nil,
			bb2:      testBinding("app1", api.LocationCommand, "simple"),
			expected: testBinding("app1", api.LocationCommand, "simple"),
		},
		{
			name:     "happy simple same",
			bb1:      testBinding("app1", api.LocationCommand, "simple"),
			bb2:      testBinding("app1", api.LocationCommand, "simple"),
			expected: testBinding("app1", api.LocationCommand, "simple"),
		},
		{
			name: "happy simple merge",
			bb1:  testBinding("app1", api.LocationPostMenu, "simple"),
			bb2:  testBinding("app1", api.LocationCommand, "simple"),
			expected: append(
				testBinding("app1", api.LocationPostMenu, "simple"),
				testBinding("app1", api.LocationCommand, "simple")...,
			),
		},
		{
			name: "happy simple 2 apps",
			bb1:  testBinding("app1", api.LocationCommand, "simple"),
			bb2:  testBinding("app2", api.LocationCommand, "simple"),
			expected: append(
				testBinding("app1", api.LocationCommand, "simple"),
				testBinding("app2", api.LocationCommand, "simple")...,
			),
		},
		{
			name: "happy 2 simple commands",
			bb1:  testBinding("app1", api.LocationCommand, "simple1"),
			bb2:  testBinding("app1", api.LocationCommand, "simple2"),
			expected: []*api.Binding{
				{
					AppID:    "app1",
					Location: "/command",
					Bindings: []*api.Binding{
						{
							AppID:    "app1",
							Location: "id-simple1",
							Hint:     "hint-simple1",
						},
						{
							AppID:    "app1",
							Location: "id-simple2",
							Hint:     "hint-simple2",
						},
					},
				},
			},
		},
		{
			name: "happy 2 apps",
			bb1: []*api.Binding{
				{
					Location: "/post_menu",
					Bindings: []*api.Binding{
						{
							AppID:       "zendesk",
							Label:       "Create zendesk ticket",
							Description: "Create ticket in zendesk",
							Call: &api.Call{
								URL: "http://localhost:4000/create",
							},
						},
					},
				},
			},
			bb2: []*api.Binding{
				{
					Location: "/post_menu",
					Bindings: []*api.Binding{
						{
							AppID:       "hello",
							Label:       "Create hello ticket",
							Description: "Create ticket in hello",
							Call: &api.Call{
								URL: "http://localhost:4000/hello",
							},
						},
					},
				},
			},
			expected: []*api.Binding{
				{
					Location: "/post_menu",
					Bindings: []*api.Binding{
						{
							AppID:       "zendesk",
							Label:       "Create zendesk ticket",
							Description: "Create ticket in zendesk",
							Call: &api.Call{
								URL: "http://localhost:4000/create",
							},
						},
						{
							AppID:       "hello",
							Label:       "Create hello ticket",
							Description: "Create ticket in hello",
							Call: &api.Call{
								URL: "http://localhost:4000/hello",
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := mergeBindings(tc.bb1, tc.bb2)
			require.Equal(t, tc.expected, out)
		})
	}
}

func TestGetBindingsGrantedLocations(t *testing.T) {
	type TC struct {
		name        string
		locations   api.Locations
		numBindings int
	}

	for _, tc := range []TC{
		{
			name: "3 locations granted",
			locations: api.Locations{
				api.Location(api.LocationChannelHeader),
				api.Location(api.LocationPostMenu),
				api.Location(api.LocationCommand),
			},
			numBindings: 3,
		},
		{
			name: "command location granted",
			locations: api.Locations{
				api.Location(api.LocationCommand),
			},
			numBindings: 1,
		},
		{
			name: "channel header location granted",
			locations: api.Locations{
				api.Location(api.LocationChannelHeader),
			},
			numBindings: 1,
		},
		{
			name: "post dropdown location granted",
			locations: api.Locations{
				api.Location(api.LocationPostMenu),
			},
			numBindings: 1,
		},
		{
			name:        "no granted locations",
			locations:   api.Locations{},
			numBindings: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bindings := []*api.Binding{
				{
					Location: api.LocationChannelHeader,
					Bindings: []*api.Binding{
						{
							Location: "send",
						},
					},
				}, {
					Location: api.LocationPostMenu,
					Bindings: []*api.Binding{
						{
							Location: "send-me",
						},
					},
				}, {
					Location: api.LocationCommand,
					Bindings: []*api.Binding{
						{
							Location: "message",
						},
					},
				},
			}

			app1 := &api.App{
				Manifest: &api.Manifest{
					AppID:              api.AppID("app1"),
					Type:               api.AppTypeBuiltin,
					RequestedLocations: tc.locations,
				},
				GrantedLocations: tc.locations,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			proxy := newTestProxyForBindings(app1, bindings, ctrl)

			cc := &api.Context{}
			out, err := proxy.GetBindings(cc)
			require.NoError(t, err)
			require.Len(t, out, tc.numBindings)
		})
	}
}

func TestGetBindingsFQL(t *testing.T) {
	initial := []*api.Binding{
		{
			Location: api.LocationChannelHeader,
			Bindings: []*api.Binding{
				{
					Location: "send",
				},
			},
		}, {
			Location: api.LocationPostMenu,
			Bindings: []*api.Binding{
				{
					Location: "send-me",
				},
				{
					Location: "send",
				},
			},
		}, {
			Location: api.LocationCommand,
			Bindings: []*api.Binding{
				{
					Location: "message",
				}, {
					Location: "message-modal",
				}, {
					Location: "manage",
					Bindings: []*api.Binding{
						{
							Location: "subscribe",
						}, {
							Location: "unsubscribe",
						},
					},
				},
			},
		},
	}

	expected := []*api.Binding{
		{
			Location: api.LocationChannelHeader,
			Bindings: []*api.Binding{
				{
					AppID:    api.AppID("app1"),
					Location: "/channel_header/send",
				},
			},
		}, {
			Location: api.LocationPostMenu,
			Bindings: []*api.Binding{
				{
					AppID:    api.AppID("app1"),
					Location: "/post_menu/send-me",
				},
				{
					AppID:    api.AppID("app1"),
					Location: "/post_menu/send",
				},
			},
		}, {
			Location: api.LocationCommand,
			Bindings: []*api.Binding{
				{
					AppID:    api.AppID("app1"),
					Location: "/command/message",
				}, {
					AppID:    api.AppID("app1"),
					Location: "/command/message-modal",
				}, {
					AppID:    api.AppID("app1"),
					Location: "/command/manage",
					Bindings: []*api.Binding{
						{
							AppID:    api.AppID("app1"),
							Location: "/command/manage/subscribe",
						}, {
							AppID:    api.AppID("app1"),
							Location: "/command/manage/unsubscribe",
						},
					},
				},
			},
		},
	}

	app := &api.App{
		Manifest: &api.Manifest{
			AppID: api.AppID("app1"),
			Type:  api.AppTypeBuiltin,
		},
		GrantedLocations: api.Locations{
			api.Location(api.LocationChannelHeader),
			api.Location(api.LocationPostMenu),
			api.Location(api.LocationCommand),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := newTestProxyForBindings(app, initial, ctrl)

	cc := &api.Context{}
	out, err := proxy.GetBindings(cc)
	require.NoError(t, err)
	require.Equal(t, expected, out)
}

func newTestProxyForBindings(app *api.App, bindings []*api.Binding, ctrl *gomock.Controller) *Proxy {
	pluginAPI := &plugintest.API{}
	pluginAPI.On("LogDebug", mock.Anything).Return(nil)
	mm := pluginapi.NewClient(pluginAPI)

	s := mock_api.NewMockStore(ctrl)
	s.EXPECT().ListApps().Return([]*api.App{app})

	up := mock_api.NewMockUpstream(ctrl)

	p := &Proxy{
		mm:    mm,
		store: s,
		builtIn: map[api.AppID]api.Upstream{
			api.AppID("app1"): up,
		},
	}

	cr := &api.CallResponse{
		Data: bindings,
	}
	bb, _ := json.Marshal(cr)
	reader := ioutil.NopCloser(bytes.NewReader(bb))
	up.EXPECT().Roundtrip(gomock.Any()).Return(reader, nil)

	return p
}
