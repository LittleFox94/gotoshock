package types

func (m Message) GomegaString() string {
	return m.DebugString()
}
