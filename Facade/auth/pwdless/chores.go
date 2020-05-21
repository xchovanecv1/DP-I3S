package pwdless

import (
	"time"

	"gitlab.com/IIIS/backend/facade/logging"
)

func (rs *Resource) choresTicker() {
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if err := rs.Store.PurgeExpiredToken(); err != nil {
				logging.Logger.WithField("chore", "purgeExpiredToken").Error(err)
			}
		}
	}()
}
