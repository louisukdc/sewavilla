package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`                // id dengan tipe uint (atau bigint sesuai kebutuhan)
	Username string `json:"username" gorm:"type:varchar(255);unique;not null"` // username unik dan tidak boleh null
	Password string `json:"password" gorm:"type:varchar(255);not null"`        // password tidak boleh null
	Email    string `json:"email" gorm:"not null;unique;column:email;size:255"`
}
