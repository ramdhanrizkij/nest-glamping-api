package pagination

import "math"

type Params struct {
	Page    int
	PerPage int
}

func (p *Params) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 10
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
}

func (p *Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func TotalPages(total int64, perPage int) int {
	return int(math.Ceil(float64(total) / float64(perPage)))
}
