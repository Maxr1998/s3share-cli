# s3share-cli

The CLI tool for [s3share](https://github.com/Maxr1998/s3share).
It is used to upload and manage files and offers an alternative to downloading through the web interface.

## Installation

s3share-cli can be easily installed through `go`:

```bash
go install github.com/maxr1998/s3share-cli@latest
```

On Arch Linux, an AUR package is available: [s3share-cli-git](https://aur.archlinux.org/packages/s3share-cli-git).

## Configuration

The CLI expect a configuration file `config.toml` in any of the following locations:

- `$XDG_CONFIG_HOME/s3share/config.toml` (usually `~/.config/s3share/config.toml`)
- `~/.config/s3share/config.toml`
- `./config.toml` (in the current directory)

A template is provided in `config_template.toml`. The following keys are required:

- `service.host`: the hostname of the site deployed to Cloudflare Pages.
- `upload.*`: these keys correspond to the `S3_*` environment variables as described in
  the [main README](https://github.com/Maxr1998/s3share?tab=readme-ov-file#configuration).
- `kv.account_id`: your Cloudflare account ID.
- `kv.api_token`: a Cloudflare Account API token with the `Account.Workers KV Storage:Edit` permission.
- `kv.namespace_id`: the ID of the Workers KV namespace to use.

## Usage

See `s3share-cli --help`.