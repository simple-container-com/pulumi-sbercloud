# pulumi-sbercloud — a Pulumi provider for Cloud.ru Advanced (ex-SberCloud)

A static Terraform bridge over
[`sbercloud-terraform/terraform-provider-sbercloud`](https://github.com/sbercloud-terraform/terraform-provider-sbercloud),
maintained for [Simple Container](https://github.com/simple-container-com/api) so that SC can
provision Cloud.ru the same way it provisions AWS, GCP and Yandex Cloud.

Upstream Terraform provider is **MPL-2.0** (it is in turn a descendant of
[`huaweicloud/terraform-provider-huaweicloud`](https://github.com/huaweicloud/terraform-provider-huaweicloud),
also MPL-2.0); attribution is preserved and the upstream source is consumed unmodified, as a Go
module dependency. The bridge code in this repository is Apache-2.0, matching its
`simple-container-com/pulumi-yandex` sibling, from which it was templated.

## Why this repo exists

1. **There is no Pulumi provider for Cloud.ru**, first- or third-party.
2. **The dynamic Terraform bridge (`pulumi package add terraform-provider …`) is unusable from
   Russia.** It needs `registry.opentofu.org` at both SDK-generation *and* plugin-download time,
   and both that registry and `registry.terraform.io` are geo-blocked there. A *static* bridge —
   this repo — resolves through `api.github.com` and the Go module proxy only, so it works from
   both sides. This is the same reasoning that produced
   [`simple-container-com/pulumi-yandex`](https://github.com/simple-container-com/pulumi-yandex).
3. **Cloud.ru Advanced is the tier with an open provider.** Cloud.ru's *other* line, Evolution,
   ships only a binary Terraform provider (`cloud-ru/evo-terraform`), which cannot be
   statically bridged — that needs either a dynamic bridge or a hand-written provider. Advanced,
   the ex-SberCloud estate, has full Go sources in a public repository, so the classic
   fork→bridge→`PluginDownloadURL` recipe applies unchanged. Do not conflate the two tiers.

## How it differs from its `pulumi-yandex` template

One structural difference, and it is the only one that matters:
**`terraform-provider-sbercloud` is a plain `terraform-plugin-sdk/v2` provider**, so this repo
uses the **plain** `tfgen.Main` / `tfbridge.Main` entrypoints, not the `pkg/pf` `MainWithMuxer`
pair. `terraform-provider-yandex` is a *muxed* server (SDKv2 alongside a plugin-framework
provider behind `tf6muxserver`) and needs the muxer at both tfgen and runtime; bridging only one
half of it silently dropped 180 of its 268 resources. Upstream here has no
`terraform-plugin-mux` and no `terraform-plugin-framework` among its direct dependencies, so
there is no second half to route to. The comments in both `provider/cmd/*/main.go` spell out the
failure modes in each direction — mismatched entrypoints are caught by neither `tfgen` nor
`go build`, only by a real `pulumi` run.

**Go SDK only**, as in the sibling: the .NET/Node/Python targets exist in the Makefile but are
not generated or published.

## Using it

```
go get github.com/simple-container-com/pulumi-sbercloud/sdk@latest
```

The plugin binary is resolved automatically from this repo's GitHub Releases — codegen bakes
`PluginDownloadURL` into the SDK's `PkgResourceDefaultOpts`, so consumers need no
`pulumi plugin install` step and no per-resource plugin option. To install it by hand:

```
pulumi plugin install resource sbercloud 0.1.0 \
  --server github://api.github.com/simple-container-com/pulumi-sbercloud
```

## Configuration

Authentication is Keystone v3 with an access-key pair. Upstream declares an `EnvDefaultFunc` for
every key, so the environment variables below work with no explicit Pulumi config:

| Config key | Env | Notes |
|---|---|---|
| `accessKey` | `SBC_ACCESS_KEY` | required |
| `secretKey` | `SBC_SECRET_KEY` | required, secret |
| `region` | `SBC_REGION_NAME` | e.g. `ru-moscow-1` |
| `projectName` | `SBC_PROJECT_NAME` | usually equals the region |
| `authUrl` | `SBC_AUTH_URL` | defaults to `https://iam.ru-moscow-1.hc.sbercloud.ru/v3` |
| `securityToken` | `SBC_SECURITY_TOKEN` | for temporary credentials |
| `userName` / `password` | `SBC_USERNAME` / `SBC_PASSWORD` | alternative to the key pair |
| `domainId` / `domainName` | `SBC_DOMAIN_ID` / `SBC_DOMAIN_NAME` | account scope |
| `enterpriseProjectId` | `SBC_ENTERPRISE_PROJECT_ID` | |
| `insecure`, `maxRetries`, `rateLimit` | `SBC_INSECURE`, `SBC_MAX_RETRIES`, `SBC_RATE_LIMIT` | |

Note `authUrl`'s default is the legacy ex-SberCloud host, which is what upstream is tested
against. If Cloud.ru Advanced moves the IAM endpoint, override `authUrl` rather than patching
the fork.

## Maintenance

`make tfgen && make build_sdks`, then tag `vX.Y.Z` — the release workflow builds the plugin for
linux/darwin × amd64/arm64 and publishes the archives the engine expects. Tag `sdk/vX.Y.Z` too,
for the Go SDK submodule. Because `RespectSchemaVersion` is on, **the release tag and
`PROVIDER_VERSION` in the Makefile must match exactly.**

Bumping upstream is a `go get github.com/sbercloud-terraform/terraform-provider-sbercloud@vX.Y.Z`
in `provider/`, then `make tfgen`. If upstream ever adopts the plugin framework, the symptom is
resources silently missing from the schema; the fix is to mirror `pulumi-yandex`'s
`MuxShimWithPF` setup in `provider/resources.go` plus the muxer entrypoints in both `cmd/` mains.
