package models

import (
	"gorm.io/gorm"
)

type User struct {
	Id       int    `json:"id" swaggerignore:"true"`
	Name     string `json:"name"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   int    `json:"-" gorm:"column:role"`
	Role     *Role  `json:"role,omitempty" gorm:"foreignKey:RoleId" swaggerignore:"true"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "user"
}

var uSortFields = []string{"id", "name", "phone", "email"}

func UsersSortFieldValid(sort string) bool {
	for _, field := range uSortFields {
		if field == sort {
			return true
		}
	}
	return false
}

func ReadUser(db *gorm.DB, email string) (*User, error) {
	var usr User
	err := db.Find(&usr, "email = ?", email).Error
	return &usr, err
}

func ReadUserRole(db *gorm.DB, email string) (*User, error) {
	var usr User
	err := db.Preload("Role").Find(&usr, "email = ?", email).Error
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

func AddUser(db *gorm.DB, newUser *User) (int64, error) {
	query := db
	newUser.RoleId = 1
	if len(newUser.Phone) <= 0 {
		query = query.Omit("Phone")
	}
	result := query.Create(&newUser)
	return result.RowsAffected, result.Error
}

func ReadAllUsers(db *gorm.DB, limit int, offset int, email string, sort string) (*[]User, error) {
	var usr []User
	query := db.
		Order("id").
		Where("role IN (SELECT id FROM role_classifier WHERE name = ?)", "user").
		Limit(limit).Offset(offset)
	if len(email) > 0 {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	err := query.Order(sort).Find(&usr).Error
	return &usr, err
}

func ReadUserById(db *gorm.DB, id int) (*User, error) {
	var usr User
	err := db.Find(&usr, id).Error
	return &usr, err
}

func RemoveUser(db *gorm.DB, userId int) (int64, error) {
	result := db.
		Where("id = ? AND role IN (SELECT id FROM role_classifier WHERE name = ?)", userId, "user").
		Delete(&User{})
	return result.RowsAffected, result.Error
}
