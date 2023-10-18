package logger

import "regexp"

type FilterFunc func(string) bool

func FilterMatchNone(string) bool { return false }

func FilterMatchAll(string) bool { return true }

type FilterType interface {
	func(string) bool | FilterFunc | *regexp.Regexp | string
}

func FilterCompile[F FilterType](f F) FilterFunc {
	switch f := any(f).(type) {
	case func(string) bool:
		return f
	case FilterFunc:
		return f
	case *regexp.Regexp:
		return f.MatchString
	case string:
		if regexp.QuoteMeta(f) == f {
			return func(s string) bool {
				return f == s
			}
		}
		return regexp.MustCompile(f).MatchString
	default:
		panic("unknown filter type")
	}
}

func Not[F FilterType](f F) FilterFunc {
	nested := FilterCompile(f)
	return func(s string) bool {
		return !nested(s)
	}
}

func And[F FilterType, G FilterType](f F, g G) FilterFunc {
	nestedF := FilterCompile(f)
	nestedG := FilterCompile(g)
	return func(s string) bool {
		return nestedF(s) && nestedG(s)
	}
}

func Or[F FilterType, G FilterType](f F, g G) FilterFunc {
	nestedF := FilterCompile(f)
	nestedG := FilterCompile(g)
	return func(s string) bool {
		return nestedF(s) || nestedG(s)
	}
}

type GenericFilter[T FilterType] struct {
	filter FilterFunc
}

func (f *GenericFilter[T]) And(g T) *GenericFilter[FilterFunc] {
	return &GenericFilter[FilterFunc]{filter: And(f.filter, g)}
}

func (f *GenericFilter[T]) Or(g T) *GenericFilter[FilterFunc] {
	return &GenericFilter[FilterFunc]{filter: Or(f.filter, g)}
}

func (f *GenericFilter[T]) Not() *GenericFilter[FilterFunc] {
	return &GenericFilter[FilterFunc]{filter: Not(f.filter)}
}

func (f *GenericFilter[T]) Done() FilterFunc {
	return f.filter
}

func NewFilter[T FilterType](f T) *GenericFilter[T] {
	return &GenericFilter[T]{filter: FilterCompile(f)}
}
