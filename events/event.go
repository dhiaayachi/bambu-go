package events

type ReportEvent interface {
	GetType() string
	IsReportEvent()
	String() string
}

type RequestEvent interface {
	GetType() string
	IsRequestEvent()
}

func NewReportEvent(eventType string) ReportEvent {
	var event ReportEvent
	switch eventType {
	case PrintType:
		event = &PrintReport{}
	case InfoType:
		event = &InfoReport{}
	case UpgrateType:
		event = &UpgradeReport{}
	default:
		return nil
	}
	return event
}
