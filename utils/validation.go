package utils

import "TravelSphere/constants"

func IsValidStatus(status string) bool {
	return constants.AllowedStatuses[status]
}
