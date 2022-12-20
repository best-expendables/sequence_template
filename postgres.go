package sequencetemplate

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	//Require for gorm
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm/clause"
)

// PostgresGenerator postgres sequence gernerator implement
type PostgresGenerator struct {
	db                    *gorm.DB
	sequences             map[string]bool
	migrateFromDBSequence bool
}

func MigrateFromDBSequence() func(*PostgresGenerator) {
	return func(gen *PostgresGenerator) {
		gen.migrateFromDBSequence = true
	}
}

// NewPostgesGenerator create new sequence generator by postgres database
func NewPostgesGenerator(db *gorm.DB, options ...func(*PostgresGenerator)) (SequenceGenerator, error) {
	gen := PostgresGenerator{db: db, sequences: make(map[string]bool)}
	if err := gen.createGaplessSequenceTable(); err != nil {
		return nil, err
	}
	for _, o := range options {
		o(&gen)
	}
	return &gen, nil
}

// Generate generate new sequence by postgres
// If length == 0 -> return original sequence from database
// Else return padded string with 0
func (gen *PostgresGenerator) Generate(seqKey string, prefix string, length int) (string, error) {
	return gen.GenerateWithStartAt(seqKey, prefix, length, 0)
}

// Generate generate new sequence by postgres
// If length == 0 -> return original sequence from database
// Else return padded string with 0
func (gen *PostgresGenerator) GenerateWithStartAt(seqKey, prefix string, length, startAt int) (string, error) {
	seq, err := gen.getSequenceFromKey(seqKey)
	if err != nil && !gen.isSequenceNotExistError(err) {
		return "", err
	} else if err != nil && gen.isSequenceNotExistError(err) {
		if err = gen.createSequence(seqKey); err != nil {
			return "", err
		}
		seq, err = gen.getSequenceFromKey(seqKey)
		if err != nil {
			return "", err
		}
	}

	strSeq := strconv.FormatInt(int64(seq+startAt), 10)

	if length == 0 {
		return fmt.Sprintf("%s%s", prefix, strSeq), nil
	}

	return fmt.Sprintf("%s%s", prefix, lefpad(strSeq, "0", length)), nil
}

func (gen *PostgresGenerator) GetCurrentSequenceFromKey(seqKey string) (int, error) {
	var seq int
	row := gen.db.Raw(fmt.Sprintf("SELECT last_value FROM %s", seqKey)).Row()
	err := row.Scan(&seq)
	if err != nil {
		return 0, err
	}

	return seq, nil
}

func lefpad(str, padStr string, length int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(padStr, length-len(str)), str)
}

func (gen *PostgresGenerator) GenerateGapless(db *gorm.DB, seqKey string, prefix string, length int) (string, error) {
	if err := gen.createGaplessSequenceRow(db, seqKey); err != nil {
		return "", err
	}
	seq, err := gen.getGaplessSequenceFromKey(db, seqKey)
	if err != nil {
		return "", err
	}
	strSeq := strconv.FormatInt(int64(seq), 10)

	if length == 0 {
		return fmt.Sprintf("%s%s", prefix, strSeq), nil
	}

	return fmt.Sprintf("%s%s", prefix, lefpad(strSeq, "0", length)), nil
}

func (gen *PostgresGenerator) getSequenceFromKey(seqKey string) (int, error) {
	var seq int
	row := gen.db.Raw(fmt.Sprintf("SELECT nextval('%s')", seqKey)).Row()
	err := row.Scan(&seq)
	if err != nil {
		return 0, err
	}

	return seq, nil
}

func (gen *PostgresGenerator) createSequence(seqey string) error {
	query := fmt.Sprintf(`CREATE SEQUENCE IF NOT EXISTS "%s"`, seqey)
	return gen.db.Exec(query).Error
}

func (gen *PostgresGenerator) isSequenceNotExistError(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "does not exist") {
		return true
	}
	return false
}

func (gen *PostgresGenerator) createGaplessSequenceTable() error {
	if err := gen.db.AutoMigrate(&GaplessSequence{}); err != nil {
		return err
	}
	return nil
}

func (gen *PostgresGenerator) createGaplessSequenceRow(db *gorm.DB, seqKey string) error {
	if gen.sequences[seqKey] {
		return nil
	}
	dbResult := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&GaplessSequence{SequenceKey: seqKey})
	if dbResult.Error != nil {
		return dbResult.Error
	}
	if dbResult.RowsAffected == 0 || !gen.migrateFromDBSequence {
		gen.sequences[seqKey] = true
		return nil
	}
	currentSequence, err := gen.GetCurrentSequenceFromKey(seqKey)
	if gen.isSequenceNotExistError(err) {
		gen.sequences[seqKey] = true
		return nil
	}
	if err != nil {
		return err
	}
	err = db.Model(&GaplessSequence{}).Where("sequence_key = ?", seqKey).Update("sequence_value", currentSequence).Error
	if err != nil {
		return err
	}
	gen.sequences[seqKey] = true
	return nil
}

func (gen *PostgresGenerator) getGaplessSequenceFromKey(db *gorm.DB, seqKey string) (int, error) {
	var sequence GaplessSequence
	dbResult := db.Model(&sequence).Clauses(clause.Returning{}).Where("sequence_key = ?", seqKey).Update("sequence_value", gorm.Expr("sequence_value + 1"))
	if err := dbResult.Error; err != nil {
		return 0, err
	}
	if dbResult.RowsAffected == 0 {
		gen.sequences[seqKey] = false
		return 0, errors.New("no row affected")
	}
	return sequence.SequenceValue, nil
}
