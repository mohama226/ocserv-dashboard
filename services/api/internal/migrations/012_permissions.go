package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"your-project/internal/models"
)


func Migration012Permissions() *gormigrate.Migration {

	return &gormigrate.Migration{

		ID: "012_permissions",

		Migrate: func(tx *gorm.DB) error {


			err := tx.AutoMigrate(
				&models.Permission{},
				&models.AdminPermission{},
			)

			if err != nil {
				return err
			}


			permissions := []models.Permission{


				// Users

				{
					Key:"users.view",
					Name:"مشاهده کاربران",
					Group:"users",
				},

				{
					Key:"users.create",
					Name:"ساخت کاربر",
					Group:"users",
				},

				{
					Key:"users.edit",
					Name:"ویرایش کاربر",
					Group:"users",
				},

				{
					Key:"users.delete",
					Name:"حذف کاربر",
					Group:"users",
				},


				{
					Key:"users.reset_password",
					Name:"تغییر رمز کاربر",
					Group:"users",
				},


				{
					Key:"users.disconnect",
					Name:"قطع اتصال کاربر",
					Group:"users",
				},



				// Service


				{
					Key:"service.view",
					Name:"مشاهده وضعیت سرویس",
					Group:"service",
				},

				{
					Key:"service.restart",
					Name:"ریستارت سرویس",
					Group:"service",
				},



				// Settings


				{
					Key:"settings.view",
					Name:"مشاهده تنظیمات",
					Group:"settings",
				},

				{
					Key:"settings.edit",
					Name:"ویرایش تنظیمات",
					Group:"settings",
				},



				// Admins


				{
					Key:"admins.view",
					Name:"مشاهده ادمین‌ها",
					Group:"admins",
				},

				{
					Key:"admins.create",
					Name:"ساخت ادمین",
					Group:"admins",
				},

				{
					Key:"admins.edit",
					Name:"ویرایش ادمین",
					Group:"admins",
				},

				{
					Key:"admins.delete",
					Name:"حذف ادمین",
					Group:"admins",
				},

				{
					Key:"admins.permissions",
					Name:"تغییر دسترسی ادمین",
					Group:"admins",
				},

			}


			for _, p := range permissions {

				var count int64

				tx.Model(&models.Permission{}).
					Where("key = ?",p.Key).
					Count(&count)


				if count == 0 {

					if err := tx.Create(&p).Error; err != nil {
						return err
					}

				}

			}


			return nil

		},


		Rollback: func(tx *gorm.DB) error {

			return tx.Migrator().
				DropTable(
					&models.AdminPermission{},
					&models.Permission{},
				)

		},

	}

}
