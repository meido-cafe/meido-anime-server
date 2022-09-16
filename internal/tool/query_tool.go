package tool

type Query struct {
	values []any
}

func NewQuery() *Query {
	return &Query{values: make([]any, 0, 4)}
}

func (p *Query) Add(data ...any) {
	for _, item := range data {
		switch item.(type) {
		case []string:
			for _, item := range item.([]string) {
				p.values = append(p.values, item)
			}
		case []int64:
			for _, item := range item.([]int64) {
				p.values = append(p.values, item)
			}
		case []int:
			for _, item := range item.([]int) {
				p.values = append(p.values, item)
			}
		default:
			p.values = append(p.values, item)
		}
	}

}

func (p *Query) Values() []any {
	return p.values
}

func (p *Query) Reset() {
	p.values = make([]any, 0, len(p.values))
}
