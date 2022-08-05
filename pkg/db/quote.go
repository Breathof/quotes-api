package db

type Quote struct {
	ID      int64 `json:"id"`
	Comment string
	Author  *Author `pg:"rel:has-one" json:"author"`
}
