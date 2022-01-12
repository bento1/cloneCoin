package p2p

type MessageKind int

const (
	// MessageNewestBlock       MessageKind = 1
	// MessageAllBlocksResqust  MessageKind = 2
	// MessageAllBlocksResponse MessageKind = 3 iota쓰면 자동으로 value랑 type을 지정해준다
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksResqust
	MessageAllBlocksResponse
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}
