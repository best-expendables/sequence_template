package sequencetemplate

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGenerateGapless(t *testing.T) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: "host=localhost user=app_user port=5432 dbname=app_db sslmode=disable password=app_pass",
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		t.Error(err)
	}
	gen, err := NewPostgesGenerator(db, MigrateFromDBSequence())
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		tx := db.Begin()
		s, err := gen.GenerateGapless(tx, "test", "test", 5)
		if err != nil {
			tx.Rollback()
			t.Error(err)
		}
		tx.Commit()
		t.Log(s)
	}
}

func TestMigration(t *testing.T) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: "host=localhost user=app_user port=5432 dbname=app_db sslmode=disable password=app_pass",
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		t.Error(err)
	}
	gen, err := NewPostgesGenerator(db, MigrateFromDBSequence())
	if err != nil {
		t.Error(err)
	}
	tx := db.Begin()
	for i := 0; i < 10; i++ {
		s, err := gen.GenerateGapless(tx, "test", "test", 5)
		if err != nil {
			tx.Rollback()
			t.Error(err)
		}
		t.Log(s)
	}
	s, err := gen.GenerateGapless(tx, "test", "test", 5)
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}
	tx.Commit()
	t.Log(s)
}
