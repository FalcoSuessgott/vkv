# MCP
`vkv` ships an [MCP](https://modelcontextprotocol.io/overview) Server for LLMs to interact with `vkv` to fetch secrets from a HashiCorp Vault engine.

So far the MCP server includes one tool (`export`) that is able to return secrets from a connected Vault server.

## Roo/Cline MCP Server Configuration
!!! note
    Unfortunately the environment variable expansion (`${env:VAULT_TOKEN}`) doesn't work within the `env` block

```json
{
  "mcpServers": {
    "vkv": {
      "type": "stdio",
      "command": "vkv",
      "args": [
        "mcp"
      ],
      "env": {
        "VAULT_TOKEN": "root",
        "VAULT_ADDR": "http://127.0.0.1:8200"
      },
    }
  }
}
```