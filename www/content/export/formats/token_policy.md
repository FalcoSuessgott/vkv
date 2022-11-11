---
title: "Token Policy Matrix"
weight: 5
---

Print the current [tokens policy capabilities](https://developer.hashicorp.com/vault/docs/commands/token/capabilities) in a matrix. 

**Requires the `update` capabilities on [`/sys/capabilities-self`](https://developer.hashicorp.com/vault/api-docs/system/capabilities-self), which is set by the default policy**

### Required flags

```bash
vkv --path <path> --format=policy
```

### Demo
<div align="center">
<br>
<img src="https://media.githubusercontent.com/media/FalcoSuessgott/vkv/master/www/static/images/policy.gif" alt="drawing" width="1000"/>
</div>