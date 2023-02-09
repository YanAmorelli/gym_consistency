package models

type User struct {
	Id              string `json:"-" gorm:"column:user_id"`
	FullName        string `json:"fullname" gorm:"column:fullname"`
	Username        string `json:"username" gorm:"column:username"`
	Email           string `json:"email" gorm:"column:email"`
	Password        string `json:"password" gorm:"column:passwd"`
	ConfirmPassword string `json:"confirmation_password" gorm:"-"`
}

type ForgetPassword struct {
	Username           string `json:"username" gorm:"-"`
	OldPassword        string `json:"old_password" gorm:"-"`
	NewPassword        string `json:"new_password" gorm:"column:passwd"`
	ConfirmNewPassword string `json:"confirmation_new_password" gorm:"-"`
}
