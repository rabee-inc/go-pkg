package pagination

import "math"

// Object ... ページネーションのオブジェクト
type Object struct {
	TotalCount  int  `json:"total_count"`
	OffsetCount int  `json:"offset_count"`
	Per         int  `json:"per"`
	Total       int  `json:"total"`
	Current     int  `json:"current"`
	Next        int  `json:"next"`
	Prev        int  `json:"prev"`
	IsFirst     bool `json:"is_first"`
	IsLast      bool `json:"is_last"`
}

// Set ... ページネーションを設定する
func (m *Object) Set(totalCount int) {
	if totalCount == 0 {
		return
	}
	m.TotalCount = totalCount
	m.OffsetCount = ((m.Current - 1) * m.Per)
	m.Total = int(math.Ceil(float64(totalCount) / float64(m.Per)))
	if m.Current < m.Total {
		m.Next = m.Current + 1
	}
	if 1 < m.Current {
		m.Prev = m.Current - 1
	}
	if m.Total <= m.Current {
		m.IsLast = true
	}
	if m.Current == 1 {
		m.IsFirst = true
	}
	return
}

// New ... ページネーションを作成する
func New(page int, per int) *Object {
	return &Object{
		Current: page,
		Per:     per,
		IsFirst: false,
		IsLast:  false,
	}
}
