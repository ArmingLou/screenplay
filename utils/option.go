package utils

type Options[T comparable] map[T]string

func (o Options[T]) GetLabels() []string {
	res := make([]string, 0)
	for _, v := range o {
		res = append(res, v)
	}
	return res
}

func (o Options[T]) GetValueByLabel(label string) T {
	for k, v := range o {
		if v == label {
			return k
		}
	}
	var def T
	return def
}
