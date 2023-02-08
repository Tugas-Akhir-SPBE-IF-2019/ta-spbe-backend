package service

type JsonManipulator interface {
	Marshal(v any) ([]byte, error)
}
