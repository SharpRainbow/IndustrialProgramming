package models

import (
	"gorm.io/gorm"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "client"
}

func ReadUser(db *gorm.DB, email string) (*User, error) {
	var usr User
	err := db.Find(&usr, "email = ?", email).Error
	return &usr, err
}

func UpdateUser(db *gorm.DB, newData *User, uid int) (int64, error) {
	updates := map[string]interface{}{}
	if newData.Name != "" {
		updates["name"] = newData.Name
	}
	if newData.Phone != "" {
		updates["phone"] = newData.Phone
	}
	if newData.Password != "" {
		updates["password"] = newData.Password
	}
	result := db.Model(&User{}).Where("id = ?", uid).Updates(updates)
	return result.RowsAffected, result.Error
}
