---
title: "Generate Vault policies using the template format"
weight: 4
---

Applying the following template-snippet to vkvs template-format we can generate vault policies for each path:


```
{{ range $path, $data := . }}
path "{{ $path }}/*" {
    capabilities = [ "create", "read" ]
}
{{ end }}
```

### Demo
<div align="center">
<br>
<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/policies.gif" alt="drawing" width="1000"/>
</div>