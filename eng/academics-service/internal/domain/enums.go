package domain

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusActive    Status = "ACTIVE"
	StatusInactive  Status = "INACTIVE"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
	StatusSuspended Status = "SUSPENDED"
)

type SubjectType string

const (
	SubjectTypeTheory    SubjectType = "THEORY"
	SubjectTypePractical SubjectType = "PRACTICAL"
	SubjectTypeLab       SubjectType = "LAB"
)

type FeeFrequency string

const (
	FeeFrequencyMonthly FeeFrequency = "MONTHLY"
	FeeFrequencyYearly  FeeFrequency = "YEARLY"
	FeeFrequencyOneTime FeeFrequency = "ONE_TIME"
)

type PaymentMode string

const (
	PaymentModeCash   PaymentMode = "CASH"
	PaymentModeCheque PaymentMode = "CHEQUE"
	PaymentModeUPI    PaymentMode = "UPI"
	PaymentModeOnline PaymentMode = "ONLINE"
)

type FeeStatus string

const (
	FeeStatusPending FeeStatus = "PENDING"
	FeeStatusPartial FeeStatus = "PARTIAL"
	FeeStatusPaid    FeeStatus = "PAID"
	FeeStatusOverdue FeeStatus = "OVERDUE"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "MONDAY"
	Tuesday   DayOfWeek = "TUESDAY"
	Wednesday DayOfWeek = "WEDNESDAY"
	Thursday  DayOfWeek = "THURSDAY"
	Friday    DayOfWeek = "FRIDAY"
	Saturday  DayOfWeek = "SATURDAY"
	Sunday    DayOfWeek = "SUNDAY"
)

type EvaluationStatus string

const (
	EvaluationStatusDraft     EvaluationStatus = "DRAFT"
	EvaluationStatusSubmitted EvaluationStatus = "SUBMITTED"
	EvaluationStatusPublished EvaluationStatus = "PUBLISHED"
)

type AttendanceStatus string

const (
	AttendancePresent AttendanceStatus = "PRESENT"
	AttendanceAbsent  AttendanceStatus = "ABSENT"
	AttendanceLate    AttendanceStatus = "LATE"
	AttendanceHalfDay AttendanceStatus = "HALF_DAY"
	AttendanceLeave   AttendanceStatus = "LEAVE"
)

type CertificateType string

const (
	CertTypeIDCard    CertificateType = "ID_CARD"
	CertTypeTC        CertificateType = "TRANSFER_CERTIFICATE"
	CertTypeBonafide  CertificateType = "BONAFIDE"
	CertTypeCharacter CertificateType = "CHARACTER"
)

type EventType string

const (
	EventTypeHoliday EventType = "HOLIDAY"
	EventTypeExam    EventType = "EXAM"
	EventTypeEvent   EventType = "EVENT"
	EventTypeMeeting EventType = "MEETING"
)
