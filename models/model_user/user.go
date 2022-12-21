package model_user

type User struct {
	FullName string `json:"fullname" gorm:"column:fullname"`
	Username string `json:"username" gorm:"column:username"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"passwd" gorm:"column:passwd"`
}
