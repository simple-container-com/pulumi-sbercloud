// Copyright 2016-2024, Pulumi Corporation.
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

package sbercloud

import (
	"path"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	sbercloud "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"

	"github.com/simple-container-com/pulumi-sbercloud/provider/pkg/version"
)

// all of the token components used below.
const (
	// This variable controls the default name of the package in the package
	// registries for nodejs and python:
	mainPkg = "sbercloud"
	// modules:
	mainMod = "index" // the sbercloud module
)

//go:embed cmd/pulumi-resource-sbercloud/bridge-metadata.json
var metadata []byte

// Provider returns additional overlaid schema and metadata associated with the provider.
func Provider() tfbridge.ProviderInfo {
	// terraform-provider-sbercloud is a plain terraform-plugin-sdk/v2 provider: its
	// main.go is a single plugin.Serve over sbercloud.Provider, with no
	// terraform-plugin-mux and no terraform-plugin-framework among its direct deps
	// (checked against v1.13.1). So a straight shimv2.NewProvider is the whole story
	// here — unlike the muxed terraform-provider-yandex this repo was templated from,
	// which needs pftfbridge.MuxShimWithPF plus the muxer entrypoints in both cmd/ mains.
	//
	// It is a descendant of terraform-provider-huaweicloud (which it imports), and
	// upstream's own schema already carries SBC_*-prefixed EnvDefaultFuncs for every
	// credential, so the Config block below only needs to make those discoverable to
	// Pulumi — it is not where the env fallback is implemented.
	prov := tfbridge.ProviderInfo{
		P: shimv2.NewProvider(sbercloud.Provider()),

		Name:        "sbercloud",
		Version:     version.Version,
		DisplayName: "SberCloud",
		Publisher:   "simple-container-com",
		LogoURL:     "",
		// Resolved by the Pulumi engine straight from this repository's GitHub Releases,
		// which is what keeps the provider usable where the Terraform/OpenTofu registries
		// are geo-blocked — the same reason pulumi-yandex exists. Codegen bakes this into
		// the generated SDK's PkgResourceDefaultOpts, so consumers need no per-resource
		// option and no `pulumi plugin install` step.
		PluginDownloadURL: "github://api.github.com/simple-container-com/pulumi-sbercloud",
		Description:       "A Pulumi package for creating and managing Cloud.ru Advanced (ex-SberCloud) resources.",
		// category/cloud tag helps with categorizing the package in the Pulumi Registry.
		Keywords:   []string{"sbercloud", "cloudru", "category/cloud"},
		License:    "Apache-2.0",
		Homepage:   "https://www.pulumi.com",
		Repository: "https://github.com/simple-container-com/pulumi-sbercloud",
		// The GitHub Org for the upstream provider. Must match the TF provider module's
		// require directive, not any replace directive.
		GitHubOrg:    "sbercloud-terraform",
		MetadataInfo: tfbridge.NewProviderMetadata(metadata),

		// Upstream already defaults each of these from its own SBC_* EnvDefaultFunc
		// (sbercloud/provider.go). Restating them here is what puts the env var in the
		// generated schema and the registry docs, so a user can see where a value came
		// from. Keep this list in sync with upstream when bumping.
		//
		// Note auth_url's upstream default is the ex-SberCloud Keystone v3 endpoint,
		// https://iam.ru-moscow-1.hc.sbercloud.ru/v3 — left alone deliberately rather
		// than repointed at a cloud.ru hostname, because that default is what the
		// provider is tested against upstream. Override it via config if Cloud.ru
		// Advanced moves the endpoint.
		Config: map[string]*tfbridge.SchemaInfo{
			"access_key": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_ACCESS_KEY"}},
			},
			"secret_key": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_SECRET_KEY"}},
			},
			"security_token": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_SECURITY_TOKEN"}},
			},
			"auth_url": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_AUTH_URL"}},
			},
			"region": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_REGION_NAME"}},
			},
			"user_name": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_USERNAME"}},
			},
			"password": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_PASSWORD"}},
			},
			"project_name": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_PROJECT_NAME"}},
			},
			"account_name": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_ACCOUNT_NAME"}},
			},
			"domain_id": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_DOMAIN_ID"}},
			},
			"domain_name": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_DOMAIN_NAME"}},
			},
			"enterprise_project_id": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_ENTERPRISE_PROJECT_ID"}},
			},
			"insecure": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_INSECURE"}},
			},
			"max_retries": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_MAX_RETRIES"}},
			},
			"rate_limit": {
				Default: &tfbridge.DefaultInfo{EnvVars: []string{"SBC_RATE_LIMIT"}},
			},
		},

		Golang: &tfbridge.GolangInfo{
			// Set where the SDK is going to be published to.
			ImportBasePath: path.Join(
				"github.com/simple-container-com/pulumi-sbercloud/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			// Opt in to all available code generation features.
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			// RespectSchemaVersion ensures the SDK is generated linking to the correct
			// version of the provider. With it on, the release tag and PROVIDER_VERSION
			// in the Makefile must match exactly.
			RespectSchemaVersion: true,
		},
	}

	// MustComputeTokens maps all resources and datasources from the upstream provider
	// into Pulumi. tokens.SingleModule puts every upstream item into the main module.
	prov.MustComputeTokens(tokens.SingleModule("sbercloud_", mainMod,
		tokens.MakeStandard(mainPkg)))

	prov.MustApplyAutoAliases()
	prov.SetAutonaming(255, "-")

	return prov
}
