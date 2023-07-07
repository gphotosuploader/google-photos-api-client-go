package gphotos

// ErrDailyQuotaExceeded is returned when the Google Photos API 'All request' per
// day quota is exceeded.
//
// See: https://developers.google.com/photos/library/guides/api-limits-quotas#general-quota-limits
type ErrDailyQuotaExceeded struct{}

func (e *ErrDailyQuotaExceeded) Error() string {
	return "daily quota exceeded"
}
