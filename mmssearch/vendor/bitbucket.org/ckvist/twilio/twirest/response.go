package twirest

// TwilioResponse holds one possible resource/response depending on type of
// request plus a Status struct.
type TwilioResponse struct {
	Accounts          *AccountsResponse          `xml:"Accounts"`
	Account           *AccountResponse           `xml:"Account"`
	Calls             *CallsResponse             `xml:"Calls"`
	Call              *CallResponse              `xml:"Call"`
	Conferences       *ConferencesResponse       `xml:"Conferences"`
	Conference        *ConferenceResponse        `xml:"Conference"`
	Exception         *ExceptionResponse         `xml:"RestException"`
	Messages          *MessagesResponse          `xml:"Messages"`
	Message           *MessageResponse           `xml:"Message"`
	Notifications     *NotificationsResponse     `xml:"Notifications"`
	Notification      *NotificationResponse      `xml:"Notification"`
	OutgoingCallerIds *OutgoingCallerIdsResponse `xml:"OutgoingCallerIds"`
	OutgoingCallerId  *OutgoingCallerIdResponse  `xml:"OutgoingCallerId"`
	Participants      *ParticipantsResponse      `xml:"Participants"`
	Participant       *ParticipantResponse       `xml:"Participant"`
	Recordings        *RecordingsResponse        `xml:"Recordings"`
	Recording         *RecordingResponse         `xml:"Recording"`
	Queues            *QueuesResponse            `xml:"Queues"`
	Queue             *QueueResponse             `xml:"Queue"`
	QueueMembers      *QueueMembersResponse      `xml:"QueueMembers"`
	QueueMember       *QueueMemberResponse       `xml:"QueueMember"`
	UsageRecords      *UsageRecordsResponse      `xml:"UsageRecords"`
	ValidationRequest *ValidationRequestResponse `xml:"ValidationRequest"`
	Status            ResponseStatus
}

type ResponseStatus struct {
	Http   int
	Twilio int
	//HttpStr  string
}

type ExceptionResponse struct {
	Code     int
	Message  string
	MoreInfo string
	Status   string
}

type Page struct {
	Page            uint   `xml:"page,attr"`
	NumPages        uint   `xml:"numpages,attr"`
	PageSize        uint   `xml:"pagesize,attr"`
	Total           uint   `xml:"total,attr"`
	Start           uint   `xml:"start,attr"`
	End             uint   `xml:"end,attr"`
	Uri             string `xml:"uri,attr"`
	FirstPageUri    string `xml:"firstpageuri,attr"`
	PreviousPageUri string `xml:"previouspageuri,attr"`
	NextPageUri     string `xml:"nextpageuri,attr"`
	LastPageUri     string `xml:"lastpageuri,attr"`
}

type AccountsResponse struct {
	Page
	Account []AccountResponse
}

type AccountResponse struct {
	Sid             string
	DateCreated     string
	DateUpdated     string
	FriendlyName    string
	Type            string
	Status          string
	AuthToken       string
	Uri             string
	OwnerAccountSid string
	SubResourceUris *AccountSubUris
}

type AccountSubUris struct {
	AvailablePhoneNumbers string
	Calls                 string
	Conferences           string
	IncomingPhoneNumbers  string
	Notifications         string
	OutgoingCallerIds     string
	Recordings            string
	Transcriptions        string
	SMSMessages           string
}

type CallsResponse struct {
	Page
	Call []CallResponse
}

type CallResponse struct {
	Sid             string
	ParentCallSid   string
	DateCreated     string
	DateUpdated     string
	AccountSid      string
	To              string
	From            string
	PhoneNumberSid  string
	Status          string
	StartTime       string
	EndTime         string
	Duration        string
	Price           string
	PriceUnit       string
	Direction       string
	AnsweredBy      string
	ForwardedFrom   string
	CallerName      string
	Uri             string
	SubResourceUris *CallSubUris
}

type CallSubUris struct {
	Notifications string
	Recordings    string
}

type ConferencesResponse struct {
	Page
	Conference []ConferenceResponse
}

type ConferenceResponse struct {
	Sid             string
	AccountSid      string
	FriendlyName    string
	Status          string
	DateCreated     string
	DateUpdated     string
	Uri             string
	SubResourceUris *ConferenceSubUris
}

type ConferenceSubUris struct {
	Participants string
}

type MessagesResponse struct {
	Page
	Message []MessageResponse
}

type MessageResponse struct {
	Sid         string
	DateCreated string
	DateUpdated string
	DateSent    string
	AccountSid  string
	To          string
	From        string
	Body        string
	NumSegments string
	Status      string
	Direction   string
	Price       string
	PriceUnit   string
	ApiVersion  string
	Uri         string
}

type NotificationsResponse struct {
	Page
	Notification []NotificationResponse
}

type NotificationResponse struct {
	Sid           string
	DateCreated   string
	DateUpdated   string
	AccountSid    string
	CallSid       string
	ApiVersion    string
	Log           string
	ErrorCode     string
	MoreInfo      string
	MessageText   string
	MessageDate   string
	RequestUrl    string
	RequestMethod string
	Uri           string
	// The fields below are only included in
	// resource from 'Notification' request
	RequestVariables string
	ResponseHeaders  string
	ResponseBody     string
}

type OutgoingCallerIdsResponse struct {
	Page
	OutgoingCallerId []OutgoingCallerIdResponse
}

type OutgoingCallerIdResponse struct {
	Sid          string
	DateCreated  string
	DateUpdated  string
	FriendlyName string
	AccountSid   string
	PhoneNumber  string
	Uri          string
}

// Response from AddOutgoingCallerId
type ValidationRequestResponse struct {
	AccountSid     string
	PhoneNumber    string
	FriendlyName   string
	ValidationCode string
	CallSid        string
}

type ParticipantsResponse struct {
	Page
	Participant []ParticipantResponse
}

type ParticipantResponse struct {
	ConferenceSid          string
	AccountSid             string
	CallSid                string
	Muted                  string
	EndConferenceOnExit    string
	StartConferenceOnEnter string
	DateCreated            string
	DateUpdated            string
	Uri                    string
}

type QueuesResponse struct {
	Page
	Queue []QueueResponse
}

type QueueResponse struct {
	Sid             string
	FriendlyName    string
	CurrentSize     string
	MaxSize         string
	AverageWaitTime string
	DateCreated     string
	DateUpdated     string
	Uri             string
}

type QueueMembersResponse struct {
	Page
	QueueMember []QueueMemberResponse
}

type QueueMemberResponse struct {
	CallSid      string
	DateEnqueued string
	WaitTime     string
	Position     string
}

type RecordingsResponse struct {
	Page
	Recording []RecordingResponse
}

type RecordingResponse struct {
	Sid         string
	DateCreated string
	DateUpdated string
	AccountSid  string
	CallSid     string
	ApiVersion  string
	Uri         string
	Duration    string
}

type UsageRecordsResponse struct {
	Page
	UsageRecord []UsageRecordResponse
}

type UsageRecordResponse struct {
	Category        string
	Description     string
	AccountSid      string
	StartDate       string
	EndDate         string
	Usage           string
	UsageUnit       string
	Count           string
	CountUnit       string
	Price           string
	PriceUnit       string
	Uri             string
	SubresourceUris *UsageRecordSubUris
}

type UsageRecordSubUris struct {
	Daily     string
	Monthly   string
	Yearly    string
	AllTime   string
	Today     string
	Yesterday string
	ThisMonth string
	LastMonth string
}
