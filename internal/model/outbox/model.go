package outbox

type OutboxItem struct {
	ID    int64
	Msg   string
	Topic string
	Key   string // записать сюда trASININGididididid
}
