package sequencetemplate

import "gorm.io/gorm"

// SequenceGenerator contract
// We may have many sequence generator, each must implement this interface
type SequenceGenerator interface {
	GenerateWithStartAt(seqKey string, prefix string, length, startAt int) (string, error)
	Generate(seqKey string, prefix string, length int) (string, error)
	GetCurrentSequenceFromKey(seqKey string) (int, error)
	GenerateGapless(db *gorm.DB, seqKey string, prefix string, length int) (string, error)
}
