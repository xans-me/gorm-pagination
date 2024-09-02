package pagination

// GroupBy applies a group by clause to the query.
func (p *Paginator) GroupBy(fields ...string) *Paginator {
	for _, field := range fields {
		p.Groups = append(p.Groups, field)
	}
	return p
}
