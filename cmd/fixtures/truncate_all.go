package fixtures

import "gorm.io/gorm"

func TruncateTables(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE likes;")
	db.Exec("TRUNCATE TABLE reports;")
	db.Exec("TRUNCATE TABLE group_drops;")
	db.Exec("TRUNCATE TABLE groups CASCADE;")
	db.Exec("TRUNCATE TABLE comment_responses CASCADE;")
	db.Exec("TRUNCATE TABLE comments CASCADE;")
	db.Exec("TRUNCATE TABLE drops CASCADE;")
	db.Exec("TRUNCATE TABLE drop_notifications CASCADE;")
	db.Exec("TRUNCATE TABLE users CASCADE;")
}
