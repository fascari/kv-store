package storage

const (
	TypeMemory Type = "memory"
	TypeRedis  Type = "redis"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func (t Type) IsValid() bool {
	switch t {
	case TypeMemory, TypeRedis:
		return true
	default:
		return false
	}
}
