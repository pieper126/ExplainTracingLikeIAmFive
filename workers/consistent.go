package workers

type ConsistentWorker struct{}

func NewConsistentWorker() *ConsistentWorker {
	return &ConsistentWorker{}
}
