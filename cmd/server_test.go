package cmd

import (
	"io"
	"net/http"
	"time"
)

func (s *VaultSuite) TestServerCommand() {
	s.Run("e2e", func() {
		secrets := map[string]interface{}{
			"sub": map[string]interface{}{
				"user": "password",
			},
			"sub2": map[string]interface{}{
				"key": false,
			},
		}

		// enable kv engine
		s.Require().NoError(vaultClient.EnableKV2Engine("export"), "enabling KV engine")

		// write secrets
		for k, secrets := range secrets {
			if m, ok := secrets.(map[string]interface{}); ok {
				s.Require().NoError(vaultClient.WriteSecrets("export", k, m))
			}
		}

		// invoke vkv server
		writer = io.Discard

		cmd := NewServerCmd()
		cmd.SetArgs([]string{"-p=export", "-P=127.0.0.1:8080"})

		var err error
		go func() {
			err = cmd.Execute()
		}()

		s.Require().NoError(err, "running server")

		// receive secrets from vkv http server
		time.Sleep(5 * time.Second)

		client := http.DefaultClient

		//nolint: noctx
		resp, err := client.Get("http://127.0.0.1:8080/export")
		s.Require().NoError(err, "requesting secrets")

		data, _ := io.ReadAll(resp.Body)

		defer resp.Body.Close()

		s.Require().Equal(http.StatusOK, resp.StatusCode, "status code")
		s.Require().Equal(`export key='false'
export user='password'
`, string(data), "response body")
	})
}
