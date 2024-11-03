# Installation

## brew
```bash
brew install falcosuessgott/tap/vkv
```
## cURL
```bash
version=$(curl https://api.github.com/repos/falcosuessgott/vkv/releases/latest -s | jq .name -r)
curl -OL "https://github.com/FalcoSuessgott/vkv/releases/download/${version}/vkv_$(uname)_$(uname -m).tar.gz"
tar xzf vkv_$(uname)_$(uname -m).tar.gz
chmod u+x vkv
./vkv version
```

## Packages
`vkv` is releases RPM- & DEB packages and Windows & MacOS Binaries.

You can find and download all artifacts in the [release](https://github.com/FalcoSuessgott/vkv/releases) section.

```bash
# Ubutu / Debian
dpkg -i vkv_<version>.deb

# RHEL / CentOS / Fedora
yum localinstall vkv_<version>.rpm

# Alpine
apk add --allow-untrusted vkv_<version>.apk

# tar.gz
tar xzf vkv_<version>.tar.gz
chmod u+x ./vkv
```

## Using `go`
```bash
go install github.com/FalcoSuessgott/vkv@latest
vkv
```

## From Sources
```bash
# requires go to be installed
git clone https://github.com/FalcoSuessgott/vkv
cd vkv
go install
```

## Docker
```bash
# ghcr.io
docker run -e VAULT_ADDR="${VAULT_ADDR}" -e VAULT_TOKEN="${VAULT_TOKEN}" ghcr.io/falcosuessgott/vkv
```
