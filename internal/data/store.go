package data

import "sort"

type Institution struct {
	RoutingNumber string
	Name          string
	RTP           bool
	FedNow        bool
}

type Store struct {
	rtp   map[string]bool
	fn    map[string]bool
	names map[string]string
}

func NewStore() *Store {
	return &Store{
		rtp:   make(map[string]bool),
		fn:    make(map[string]bool),
		names: make(map[string]string),
	}
}

func (s *Store) LoadRTP(rtns []string) {
	for _, rtn := range rtns {
		s.rtp[rtn] = true
	}
}

func (s *Store) LoadFedNow(rtns []string) {
	for _, rtn := range rtns {
		s.fn[rtn] = true
	}
}

func (s *Store) LoadNames(names map[string]string) {
	for rtn, name := range names {
		s.names[rtn] = name
	}
}

func (s *Store) Lookup(rtn string) Institution {
	return Institution{
		RoutingNumber: rtn,
		Name:          s.names[rtn],
		RTP:           s.rtp[rtn],
		FedNow:        s.fn[rtn],
	}
}

func (s *Store) All() []Institution {
	seen := make(map[string]bool)
	for rtn := range s.rtp {
		seen[rtn] = true
	}
	for rtn := range s.fn {
		seen[rtn] = true
	}

	result := make([]Institution, 0, len(seen))
	for rtn := range seen {
		result = append(result, s.Lookup(rtn))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RoutingNumber < result[j].RoutingNumber
	})
	return result
}
