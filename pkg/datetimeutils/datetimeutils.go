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

func RFC3339ToEpoch(rfc3339 string) int64 {
	time, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		return 0
	}

	return time.Unix()
}

func IsValidRFC3339(rfc3339 string) bool {
	_, err := time.Parse(time.RFC3339, rfc3339)
	return err == nil
}
