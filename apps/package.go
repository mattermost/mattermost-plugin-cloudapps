// Copyright (c) 2021-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

// Package apps provides the data types, constants, and convenience functions
// for the Mattermost Apps API.
//
// Function
//
// Functions are invoked in response to a user or a notification event. A
// Function to invoke is described by a Call, and is passed a CallRequest when
// invoked. For user-invoked functions, the inputs can be collected from the
// user with a Form, either as a modal, or as a /command with autocomplete.
//
// Call
//
// A Call
// (https://pkg.go.dev/github.com/mattermost/mattermost-plugin-apps/apps#Call)
// is a general request to an App server. Calls are used to fetch App's
// bindings, and to process user input, webhook events, and dynamic data
// lookups. The Call struct defines what function to invoke, and how.
//
// CallRequest
// (https://pkg.go.dev/github.com/mattermost/mattermost-plugin-apps/apps#CallRequest)
// is passed to functions at invocation time. It includes the originating Call,
// and adds user-input values, Context, etc. CallResponse
// (https://pkg.go.dev/github.com/mattermost/mattermost-plugin-apps/apps#CallResponse)
// is the result of a call.
//
// Context and Expand
//
// Context
// (https://pkg.go.dev/github.com/mattermost/mattermost-plugin-apps/apps#Context)
// of a CallRequest sent to the App includes the IDs of relevant Mattermost
// entities, such as the Location the call originated from, acting user ID,
// channel ID, etc. It also includes the Mattermost site URL, and the access
// tokens requested by the Call. This allows the function to use the Mattermost
// REST API.
//
// Context may be further expanded using the Expand attribute of a call to
// include the detailed data on the User, Channel, Team, and other relevant
// entities. Expanded context may also contain 3rd party OAuth2 access token
// previously maintained and stored in Mattermost.
//
// Expand
// (https://pkg.go.dev/github.com/mattermost/mattermost-plugin-apps/apps#Expand)
// specifies what entities should be included in the call's expanded context,
// and how much data should be included. Expanding a call requests result in
// ExpandedContext
// (https://pkg.go.dev/github.com/mattermost/mattermost-plugin-apps/apps#ExpandedContext)
// filled in.
//
// Special Notes
//
// ### Use of router packages in Apps - Go (gorilla mux) - JavaScript
//
// ### Call vs Notification
//
// ### AWS Lambda packaging
package apps
