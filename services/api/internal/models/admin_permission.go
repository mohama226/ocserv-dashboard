package models

type AdminPermission struct {
	ID uint `gorm:"primaryKey"`

	AdminID uint
	PermissionID uint

	Admin User `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE"`
	Permission Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}
