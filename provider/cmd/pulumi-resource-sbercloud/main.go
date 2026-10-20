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
	_ "embed"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"

	sbercloud "github.com/simple-container-com/pulumi-sbercloud/provider"
	"github.com/simple-container-com/pulumi-sbercloud/provider/pkg/version"
)

//go:embed schema.json
var pulumiSchema []byte

func main() {
	// The plain tfbridge.Main, pairing with the plain tfgen.Main in
	// cmd/pulumi-tfgen-sbercloud. See the long comment there for why this repo does
	// NOT use the pf/MainWithMuxer entrypoints its pulumi-yandex template does:
	// terraform-provider-sbercloud is terraform-plugin-sdk/v2 only, so there is no
	// second half to route to.
	//
	// The two must agree. Mismatched entrypoints are not caught by tfgen or by
	// `go build` — the yandex sibling learned this the hard way, where the plain Main
	// over a muxed shim panicked on the engine's first runtime call with
	// "schemaOnlyProvider does not implement runtime operation InitLogging".
	tfbridge.Main("sbercloud", version.Version, sbercloud.Provider(), pulumiSchema)
}
