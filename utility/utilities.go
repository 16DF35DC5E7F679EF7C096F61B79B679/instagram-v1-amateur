package utility

import "time"

func GetHumanReadableTime(epoch int64) string {
	return time.Unix(epoch, 0).Format(time.ANSIC)
}
