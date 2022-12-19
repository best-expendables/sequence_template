package sequencetemplate

type GaplessSequence struct {
	ID            uint   `gorm:"primaryKey"`
	SequenceKey   string `gorm:"uniqueIndex"`
	SequenceValue int
}
