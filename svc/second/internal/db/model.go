package db

import "time"

const (
	Table   = "`model`"
	Columns = "`id`, `name`, `updated_at`"
)

type Model struct {
	ID        int64
	Name      string
	UpdatedAt time.Time
}

func (m *Model) ToPtrArgs() []any {
	return []any{
		&m.ID, &m.Name, &m.UpdatedAt,
	}
}
