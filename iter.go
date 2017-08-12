package icheck

// Query is the function used to get a page listing.
type Query func(*RequestValues) ([]interface{}, error)

// Iter provides a convenient interface
// for iterating over the elements
// returned from paginated list API calls.
// Successive calls to the Next method
// will step through each item in the list,
// fetching pages of items as needed.
// Iterators are not thread-safe, so they should not be consumed
// across multiple goroutines.
type Iter struct {
	query  Query
	qs     *RequestValues
	values []interface{}
	params ListParams
	err    error
	cur    interface{}
}

// GetIter returns a new Iter for a given query and its options.
func GetIter(params *ListParams, qs *RequestValues, query Query) *Iter {
	iter := &Iter{}
	iter.query = query

	p := params
	if p == nil {
		p = &ListParams{}
	}
	iter.params = *p

	q := qs
	if q == nil {
		q = &RequestValues{}
	}
	iter.qs = q

	return iter
}
