// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfgen"

	sbercloud "github.com/simple-container-com/pulumi-sbercloud/provider"
	"github.com/simple-container-com/pulumi-sbercloud/provider/pkg/version"
)

func main() {
	// The plain tfgen.Main, NOT pf/tfgen.MainWithMuxer — and this is the one place
	// this repo deliberately diverges from its sibling simple-container-com/pulumi-yandex,
	// which was its template.
	//
	// terraform-provider-yandex is a *muxed* server (terraform-plugin-sdk/v2 alongside a
	// terraform-plugin-framework provider behind tf6muxserver), so bridging it needs the
	// muxer entrypoints at both tfgen and runtime. terraform-provider-sbercloud is not:
	// its main.go is a single `plugin.Serve(&plugin.ServeOpts{ProviderFunc: sbercloud.Provider})`
	// over terraform-plugin-sdk/v2 only, with no terraform-plugin-mux and no
	// terraform-plugin-framework among its direct dependencies (checked against v1.13.1).
	//
	// Using MainWithMuxer here anyway would not fail loudly: it would write a one-entry
	// mux dispatch table into bridge-metadata.json and build fine. Keep the plain Main
	// unless upstream actually adopts the plugin framework — if it ever does, the symptom
	// is resources silently missing from the schema, and the fix is to mirror
	// pulumi-yandex's MuxShimWithPF setup in resources.go plus the muxer entrypoints in
	// both cmd/ mains.
	tfgen.Main("sbercloud", version.Version, sbercloud.Provider())
}
