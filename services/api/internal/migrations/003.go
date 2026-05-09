package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/mmtaee/ocserv-dashboard/common/pkg/logger"
	"gorm.io/gorm"
)

var Migration003 = &gormigrate.Migration{
	ID: "003_fix_traffic_statistics_ocserv_user_fk",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Exec(`
			DELETE FROM ocserv_user_traffic_statistics t
			WHERE NOT EXISTS (
				SELECT 1 FROM ocserv_users u WHERE u.id = t.oc_user_id
			);
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE ocserv_user_traffic_statistics
			DROP CONSTRAINT IF EXISTS fk_traffic_user;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE ocserv_user_traffic_statistics
			DROP CONSTRAINT IF EXISTS fk_traffic_ocserv_user;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE ocserv_user_traffic_statistics
			ADD CONSTRAINT fk_traffic_ocserv_user
			FOREIGN KEY (oc_user_id)
			REFERENCES ocserv_users(id)
			ON DELETE CASCADE;
		`).Error; err != nil {
			return err
		}

		logger.Info("migration 003 fixed traffic statistics ocserv user foreign key successfully")
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE ocserv_user_traffic_statistics
			DROP CONSTRAINT IF EXISTS fk_traffic_ocserv_user;
		`).Error
	},
}
