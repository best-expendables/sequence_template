package sequencetemplate

// SequenceGenerator contract
// We may have many sequence generator, each must implement this interface
type SequenceGenerator interface {
	GenerateWithStartAt(seqKey string, prefix string, length, startAt int) (string, error)
	Generate(seqKey string, prefix string, length int) (string, error)
}
