package vault

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

const (
	renewalInterval  = 10 * time.Second // default renewal interval
	renewalIncrement = 30
)

// LeaseRefresher periodically checks the ttl of the current lease and attempts to renew it if the ttl is less than half of the creation ttl.
// if the token renewal fails, a new login with the configured auth method is performed
// this func is supposed to run as a goroutine.
// nolint: funlen, gocognit, cyclop
func (v *Vault) LeaseRefresher(ctx context.Context) {
	ticker := time.NewTicker(renewalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			token, err := v.Client.Auth().Token().LookupSelf()
			if err != nil {
				slog.Error("failed to lookup token", slog.Any("error", err))
				continue
			}

			creationTTL, ok := token.Data["creation_ttl"].(json.Number)
			if !ok {
				slog.Error("failed to assert creation_ttl type")

				continue
			}

			ttl, ok := token.Data["ttl"].(json.Number)
			if !ok {
				slog.Error("failed to assert ttl type")

				continue
			}

			creationTTLFloat, err := creationTTL.Float64()
			if err != nil {
				slog.Error("failed to parse creation_ttl", slog.Any("error", err))

				continue
			}

			ttlFloat, err := ttl.Float64()
			if err != nil {
				slog.Error("failed to parse ttl", slog.Any("error ", err))

				continue
			}

			slog.Info("checking token renewal", slog.Float64("creation_ttl", creationTTLFloat), slog.Float64("ttl", ttlFloat))

			//nolint: nestif
			if ttlFloat < creationTTLFloat/2 {
				slog.Info("attempting token renewal", slog.Duration("renewal_seconds", renewalInterval))

				if _, err := v.Client.Auth().Token().RenewSelf(renewalIncrement); err != nil {
					slog.Error("failed to renew token, performing new authentication", slog.Any("error", err))
				} else {
					slog.Info("successfully refreshed token")
				}
			} else {
				slog.Info("skipping token renewal")
			}

		case <-ctx.Done():
			slog.Info("token refresher shutting down")

			return
		}
	}
}
