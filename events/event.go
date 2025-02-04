package events

type ReportEvent interface {
	GetType() string
	IsReportEvent()
	String() string
}

type RequestEvent interface {
	IsRequestEvent()
	String() string
	SetSeq(id uint64)
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
	case SystemType:
		event = &LedCtlReport{}
	default:
		return nil
	}
	return event
}
