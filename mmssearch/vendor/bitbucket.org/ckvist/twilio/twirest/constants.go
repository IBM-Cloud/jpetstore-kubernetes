package twirest

// Call, Message & Conference status strings
const (
	TwiInit       = "init"
	TwiQueued     = "queued"
	TwiSending    = "sending"
	TwiSent       = "sent"
	TwiReceiving  = "receiving"
	TwiReceived   = "received"
	TwiFailed     = "failed"
	TwiRinging    = "ringing"
	TwiInProgress = "in-progress"
	TwiCanceled   = "canceled"
	TwiCompleted  = "completed"
	TwiBusy       = "busy"
	TwiNoAnswer   = "no-answer"
)

// UsageRecords categories
const (
	TwiCalls                   = "calls"
	TwiCallsInbound            = "calls-inbound"
	TwiCallsInboundLocal       = "calls-inbound-local"
	TwiCallsInboundTollfree    = "calls-inbound-tollfree"
	TwiCallsOutbound           = "calls-outbound"
	TwiCallsClient             = "calls-client"
	TwiCallsSip                = "calls-sip"
	TwiSms                     = "sms"
	TwiSmsInbound              = "sms-inbound"
	TwiSmsInboundShortcode     = "sms-inbound-shortcode"
	TwiSmsInboundLongcode      = "sms-inbound-longcode"
	TwiPhoneNumbers            = "phonenumbers"
	TwiPhoneNumbersTollFree    = "phonenumbers-tollfree"
	TwiPhoneNumbersLocal       = "phonenumbers-local"
	TwiShortcodes              = "shortcodes"
	TwiShortcodesVanity        = "shortcodes-vanity"
	TwiShortcodesRandom        = "shortcodes-random"
	TwiShortcodesCustomerOwned = "shortcodes-customerowned"
	TwiCallerIdLookups         = "calleridlookups"
	TwiRecordings              = "recordings"
	TwiTranscriptions          = "transcriptions"
	TwiRecordingStorage        = "recordingstorage"
	TwiTotalPrice              = "totalprice"
)

// UsageRecords subresources
const (
	TwiDaily     = "Daily"
	TwiMonthly   = "Monthly"
	TwiYearly    = "Yearly"
	TwiAllTime   = "AllTime"
	TwiToday     = "Today"
	TwiYesterday = "Yesterday"
	TwiThisMonth = "ThisMonth"
	TwiLastMonth = "LastMonth"
)

// Account status strings
const (
	TwiClosed    = "closed"
	TwiSuspended = "suspended"
	TwiActive    = "active"
)
