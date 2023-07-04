package shared

const (
	Mobile DeviceType = iota
	Tablet
	Desktop
)
const (
	Instant NotificationType = iota
	Scheduled
)
const (
	Sms Channel = iota
	Email
	Sms_Email
)
const (
	Pending NotificationStatus = iota
	Sent
	Received
	Failed
)
