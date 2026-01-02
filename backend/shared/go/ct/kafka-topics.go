package ct

type KafkaTopic struct {
	index int
	topic string
}

var (
	Notification = KafkaTopic{1, "notification"}
)
