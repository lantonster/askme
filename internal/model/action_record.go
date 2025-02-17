package model

type ActionRecordType string

const (
	ActionRecordTypeEmail            = "email"
	ActionRecordTypePassword         = "password"
	ActionRecordTypeEditUserinfo     = "edit_userinfo"
	ActionRecordTypeQuestion         = "question"
	ActionRecordTypeAnswer           = "answer"
	ActionRecordTypeComment          = "comment"
	ActionRecordTypeEdit             = "edit"
	ActionRecordTypeInvitationAnswer = "invitation_answer"
	ActionRecordTypeSearch           = "search"
	ActionRecordTypeReport           = "report"
	ActionRecordTypeDelete           = "delete"
	ActionRecordTypeVote             = "vote"
)

// ActionRecord [Redis] 行为记录
type ActionRecord struct {
	Num      int    `json:"num"`
	LastTime int64  `json:"last_time"`
	Config   string `json:"config"`
}
