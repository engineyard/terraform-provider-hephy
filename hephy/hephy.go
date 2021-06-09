package hephy

import (
	hephy "github.com/teamhephy/controller-sdk-go"
)

func CheckConnection(client *hephy.Client) error {
	err := client.CheckConnection()

	if err != nil {
		// We don't care about controller API version mismatches
		if err != hephy.ErrAPIMismatch {
			return err
		}
	}

	return nil
}
