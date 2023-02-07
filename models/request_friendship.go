package models

type RequestFriendship struct {
	RequestId          int    `json:"request_id" gorm:"column:request_id"`
	UserSent           string `json:"user_sent" gorm:"column:user_sent"`
	UserReceived       string `json:"user_received" gorm:"column:user_received"`
	DateSentRequest    string `json:"date_sent_request" gorm:"column:dt_sented"`
	RepliedRequestDate string `json:"date_replied_request" gorm:"column:dt_replied"`
	RequestStatus      string `json:"request_status" gorm:"column:request_status"`
}
