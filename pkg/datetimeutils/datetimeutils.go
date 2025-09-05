package datetimeutils

import "time"

func EpochNow() int64 {
	return time.Now().Unix()
}

func EpochToTime(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

func TimeToEpoch(time time.Time) int64 {
	return time.Unix()
}

func EpochToRFC3339(epoch int64) string {
	return EpochToTime(epoch).Format(time.RFC3339)
}
