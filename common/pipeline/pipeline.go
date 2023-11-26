package pipeline

import "context"

//type Filter[T any] interface {
//	Process(ctx context.Context, object T) error
//}
//
//type PipelineController[T any] struct {
//	name      string
//	filters   *[]Filter[T]
//	needBreak bool
//}
//
//func NewPipelineController[T any](name string, filters ...Filter[T]) *PipelineController[T] {
//	return &PipelineController[T]{
//		name:    name,
//		filters: &filters,
//	}
//}
//
//func (f *PipelineController[T]) Process(ctx context.Context, object T) error {
//	for _, filter := range *f.filters {
//		if f.needBreak {
//			return nil
//		}
//		err := filter.Process(ctx, object)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (f *PipelineController[T]) Name() string {
//	return f.name
//}

type HandlerFunc[T any] func(ctx context.Context, object T) error

type Filter[T any] func(next HandlerFunc[T]) HandlerFunc[T]

func FilterChain[T any](filters ...Filter[T]) Filter[T] {
	return func(next HandlerFunc[T]) HandlerFunc[T] {
		for i := len(filters) - 1; i >= 0; i-- {
			next = filters[i](next)
		}
		return next
	}
}

func (f Filter[T]) Then(h HandlerFunc[T]) HandlerFunc[T] {
	return f(h)
}
