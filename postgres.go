package sequencetemplate

import (
	"fmt"
	"github.com/jinzhu/gorm"
	//Require for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

// PostgresGenerator postgres sequence gernerator implement
type PostgresGenerator struct {
	db *gorm.DB
}

// NewPostgesGenerator create new sequence generator by postgres database
// Need postges connection config set in env
// POSTGRES_HOST
// POSTGRES_PORT
// POSTGRES_DB
// POSTGRES_USER
// POSTGRES_PASS
func NewPostgesGenerator() (SequenceGenerator, error) {
	dbConf := getAppConfigFromEnv()
	return NewPostgesGeneratorFromConfig(dbConf)
}

// NewPostgesGeneratorFromConfig create generator by passing config parameter
func NewPostgesGeneratorFromConfig(dbConf *PostgresConfig) (SequenceGenerator, error) {
	conn, err := gorm.Open("postgres", dbConf.getConnectionURL())
	if err != nil {
		return nil, err
	}
	generator := PostgresGenerator{db: conn}

	return &generator, nil
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
	seq, err := gen.GetSequenceFromKey(seqKey)
	if err != nil && !gen.isSequenceNotExistError(err) {
		return "", err
	} else if err != nil && gen.isSequenceNotExistError(err) {
		if err = gen.createSequence(seqKey); err != nil {
			return "", err
		}
		seq, err = gen.GetSequenceFromKey(seqKey)
		if err != nil {
			return "", err
		}
	}

	strSeq := strconv.FormatInt(int64(seq + startAt), 10)

	if length == 0 {
		return fmt.Sprintf("%s%s", prefix, strSeq), nil
	}

	return fmt.Sprintf("%s%s", prefix, lefpad(strSeq, "0", length)), nil
}

func lefpad(str, padStr string, length int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(padStr, length-len(str)), str)
}

func (gen *PostgresGenerator) GetSequenceFromKey(seqKey string) (int, error) {
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
	if pqErr, ok := (err.(*pq.Error)); ok {
		if pqErr.Code == pq.ErrorCode("42P01") {
			return true
		}
	}

	return false
}
