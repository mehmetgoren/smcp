package abstract

type Repository[T any] interface {
	Get(id string) (T, error)
}
