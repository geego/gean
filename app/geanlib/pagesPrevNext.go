package geanlib

// Prev returns the previous page reletive to the given page.
func (p Pages) Prev(cur *Page) *Page {
	for x, c := range p {
		if c.UniqueID() == cur.UniqueID() {
			if x == 0 {
				return p[len(p)-1]
			}
			return p[x-1]
		}
	}
	return nil
}

// Next returns the next page reletive to the given page.
func (p Pages) Next(cur *Page) *Page {
	for x, c := range p {
		if c.UniqueID() == cur.UniqueID() {
			if x < len(p)-1 {
				return p[x+1]
			}
			return p[0]
		}
	}
	return nil
}
