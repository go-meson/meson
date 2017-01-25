package object

type ObjectType int

type ObjectRef interface {
	GetID() int64
	GetObjectType() ObjectType
}
