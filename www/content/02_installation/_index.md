---
title: "Installation"
weight: 2
---

### cURL
Download vkv using curl for Linux x86_64 machines:
```bash
curl -0L https://github.com/FalcoSuessgott/vkv/releases/latest/download/vkv_0.2.0_$(uname)_$(uname -m).tar.gz
tar xzf vkv_0.2.0_Linux_x86_64.tar.gz
chmod u+x vkv
./vkv version
vkv 0.2.0
```

### Packages
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

### Using `go`
```bash
go install github.com/FalcoSuessgott/vkv@latest
vkv
```

### From Sources
```bash
# requires go to be installed
git clone https://github.com/FalcoSuessgott/vkv
cd vkv
go install 
```
