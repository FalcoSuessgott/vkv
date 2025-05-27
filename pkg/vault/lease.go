package vault

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"
)

const (
	defaultRenewalInterval  = 10 * time.Second
	defaultRenewalIncrement = 30

	envVarRenewalInterval  = "VKV_RENEWAL_INTERVAL"
	envVarRenewalIncrement = "VKV_RENEWAL_INCREMENT"
	envVarRefresherEnabled = "VKV_LEASE_REFRESHER_ENABLED"
)

// LeaseRefresher periodically checks the ttl of the current lease and attempts to renew it if the ttl is less than half of the creation ttl.
// if the token renewal fails, a new login with the configured auth method is performed
// this func is supposed to run as a goroutine.
// nolint: funlen, gocognit, cyclop
func (v *Vault) LeaseRefresher(ctx context.Context) {
	if _, ok := os.LookupEnv(envVarRefresherEnabled); ok {
		return
	}

	renewalInterval := defaultRenewalInterval

	if v, ok := os.LookupEnv(envVarRenewalInterval); ok {
		if i, err := strconv.Atoi(v); err == nil {
			renewalInterval = time.Duration(i) * time.Second
		}
	}

	renewalIncrement := defaultRenewalIncrement

	if v, ok := os.LookupEnv(envVarRenewalIncrement); ok {
		if i, err := strconv.Atoi(v); err == nil {
			renewalIncrement = i
		}
	}

	ticker := time.NewTicker(renewalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			token, err := v.Client.Auth().Token().LookupSelfWithContext(ctx)
			if err != nil {
				continue
			}

			creationTTL, ok := token.Data["creation_ttl"].(json.Number)
			if !ok {
				continue
			}

			ttl, ok := token.Data["ttl"].(json.Number)
			if !ok {
				continue
			}

			creationTTLFloat, err := creationTTL.Float64()
			if err != nil {
				continue
			}

			ttlFloat, err := ttl.Float64()
			if err != nil {
				continue
			}

			//nolint: nestif
			if ttlFloat < creationTTLFloat/2 {
				if _, err := v.Client.Auth().Token().RenewSelfWithContext(ctx, renewalIncrement); err != nil {
					continue
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
