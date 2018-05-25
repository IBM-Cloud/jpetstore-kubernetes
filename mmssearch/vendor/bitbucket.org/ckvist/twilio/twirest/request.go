package twirest

// uri URI resource
// Used for the request resource, NOTE: only the tag is used
type uri struct {
}

// Request a list of the account resources
type Accounts struct {
	FriendlyName string `FriendlyName=`
	Status       string `Status=`
}

// Account resource information for a single account
type Account struct {
	Sid string
}

// Request list of calls made to and from account
type Calls struct {
	resource        uri    `/Calls`
	To              string `To=`
	From            string `From=`
	Status          string `Status=`
	StartTime       string `StartTime=`
	StartTimeBefore string `StartTime<=`
	StartTimeAfter  string `StartTime>=`
	ParentCallSid   string `ParentCallSid=`
}

// Request call information about a single call
type Call struct {
	resource      uri    `/Calls`
	Sid           string // CallSid
	Recordings    bool
	Notifications bool
}

// Request to make a phone call
type MakeCall struct {
	resource             uri    `/Calls`
	From                 string `From=`
	To                   string `To=`
	Url                  string `Url=`
	ApplicationSid       string `ApplicationSid=`
	Method               string `Method=`
	FallbackUrl          string `FallbackUrl=`
	FallbackMethod       string `FallbackMethod=`
	StatusCallback       string `StatusCallback=`
	StatusCallbackMethod string `StatusCallbackMethod=`
	SendDigits           string `SendDigits=`
	IfMachine            string `IfMachine=`
	Timeout              string `Timeout=`
	Record               string `Record=`
	SipAuthUsername      string `SipAuthUsername=`
	SipAuthPassword      string `SipAuthPassword=`
}

// Request to modify call in queue/progress
type ModifyCall struct {
	resource             uri `/Calls`
	Sid                  string
	Url                  string `Url=`
	Method               string `Method=`
	Status               string `Status=`
	FallbackUrl          string `FallbackUrl=`
	FallbackMethod       string `FallbackMethod=`
	StatusCallback       string `StatusCallback=`
	StatusCallbackMethod string `StatusCallbackMethod=`
}

// List conferences within an account
type Conferences struct {
	resource          uri    `/Conferences`
	Status            string `Status=`
	FriendlyName      string `FriendlyName=`
	DateCreated       string `DateCreated=`
	DateCreatedBefore string `DateCreated<=`
	DateCreatedAfter  string `DateCreated>=`
	DateUpdated       string `DateUpdated=`
	DateUpdatedBefore string `DateUpdated<=`
	DateUpdatedAfter  string `DateUpdated>=`
}

// Resource for individual conference instance
type Conference struct {
	resource uri `/Conferences`
	Sid      string
}

// Request list of participants in a conference
type Participants struct {
	resource    uri    `/Conferences`
	subresource uri    `/Participants`
	Sid         string // Conference Sid
	Muted       string `Muted=`
}

// Resource about single conference participant
type Participant struct {
	resource    uri    `/Conferences`
	subresource uri    `/Participants`
	Sid         string // Conference Sid
	CallSid     string // required field
}

// Remove a participant from a conference
type DeleteParticipant struct {
	resource    uri    `/Conferences`
	subresource uri    `/Participants`
	Sid         string // Conference Sid
	CallSid     string // required field
}

// Request to change the status of a participant
type UpdateParticipant struct {
	resource    uri    `/Conferences`
	subresource uri    `/Participants`
	Sid         string // Conference Sid
	CallSid     string // required field
	Muted       string `Muted=`
}

// Messages struct for request of list of messages
type Messages struct {
	resource       uri    `/Messages`
	To             string `To=`
	From           string `From=`
	DateSent       string `DateSent=`
	DateSentBefore string `DateSent<=`
	DateSentAfter  string `DateSent>=`
}

// Message struct for request of single message
type Message struct {
	resource uri    `/Messages`
	Sid      string // MessageSid
	Media    bool
	MediaSid string
}

// Message struct for request to send a message
type SendMessage struct {
	resource            uri    `/Messages`
	Text                string `Body=`
	MediaUrl            string `MediaUrl=`
	From                string `From=`
	To                  string `To=`
	MessagingServiceSid string `MessagingServiceSid=`
	ApplicationSid      string `ApplicationSid=`
	StatusCallback      string `StatusCallback=`
}

// Notifications struct for request of a possible list of notifications
type Notifications struct {
	resource      uri    `/Notifications`
	Log           string `Log=`
	MsgDate       string `MessageDate=`
	MsgDateBefore string `MessageDate<=`
	MsgDateAfter  string `MessageDate>=`
}

// Notification struct for request of a specific notification
type Notification struct {
	resource uri `/Notifications`
	Sid      string
}

// DeleteNotification struct for removal of a notification
type DeleteNotification struct {
	resource uri `/Notifications`
	Sid      string
}

// Get outgoing caller IDs
type OutgoingCallerIds struct {
	resource     uri    `/OutgoingCallerIds`
	PhoneNumber  string `PhoneNumber=`
	FriendlyName string `FriendlyName=`
}

// Get outgoing caller ID
type OutgoingCallerId struct {
	resource uri `/OutgoingCallerIds`
	Sid      string
}

type UpdateOutgoingCallerId struct {
	resource     uri `/OutgoingCallerIds`
	Sid          string
	FriendlyName string `FriendlyName=`
}

type DeleteOutgoingCallerId struct {
	resource uri `/OutgoingCallerIds`
	Sid      string
}

type AddOutgoingCallerId struct {
	resource             uri    `/OutgoingCallerIds`
	PhoneNumber          string `PhoneNumber=`
	FriendlyName         string `FriendlyName=`
	CallDelay            string `CallDelay=`
	Extension            string `Extension=`
	StatusCallback       string `StatusCallback=`
	StatusCallbackMethod string `StatusCallbackMethod=`
}

// List recordings resource
type Recordings struct {
	resource          uri    `/Recordings`
	CallSid           string `CallSid=`
	DateCreated       string `DateCreated=`
	DateCreatedBefore string `DateCreated<=`
	DateCreatedAfter  string `DateCreated>=`
}

// Request resource for an individual recording
type Recording struct {
	resource uri    `/Recordings`
	Sid      string // RecordingSid
}

// Delete a recording
type DeleteRecording struct {
	resource uri    `/Recordings`
	Sid      string // RecordingSid
}

// Request usage by the account
type UsageRecords struct {
	resource    uri `/Usage/Records`
	SubResource string
	Category    string `Category=`
	StartDate   string `StartDate=`
	EndDate     string `EndDate=`
}

// List queues within an account
type Queues struct {
	resource uri `/Queues`
}

// Get resource for an individual Queue instance
type Queue struct {
	resource uri    `/Queues`
	Sid      string // QueueSid
}

// Create a new queue
type CreateQueue struct {
	resource     uri    `/Queues`
	FriendlyName string `FriendlyName=`
	MaxSize      string `MaxSize=`
}

// Request to change queue properties
type ChangeQueue struct {
	resource     uri `/Queues`
	Sid          string
	FriendlyName string `FriendlyName=`
	MaxSize      string `MaxSize=`
}

// Remove a queue
type DeleteQueue struct {
	resource uri    `/Queues`
	Sid      string // QueueSid
}

// List members of a queue
type QueueMembers struct {
	resource    uri    `/Queues`
	subresource uri    `/Members`
	Sid         string // QueueSid
}

// Request resource for a queue member
type QueueMember struct {
	resource    uri    `/Queues`
	subresource uri    `/Members`
	Sid         string // QueueSid
	CallSid     string // either this field or Front is required
	Front       bool
}

// Remove a member from a queue and redirect the member's call to a TwiML site
type DeQueue struct {
	resource    uri    `/Queues`
	subresource uri    `/Members`
	Sid         string // Queue Sid
	CallSid     string // either this field or Front is required
	Front       bool
	Url         string `Url=`
	Method      string `Method=`
}
