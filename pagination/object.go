package pagination

import "math"

// Object ... ページネーションのオブジェクト
type Object struct {
	TotalCount  int  `json:"total_count"`
	OffsetCount int  `json:"offset_count"`
	PerPage     int  `json:"per_page"`
	TotalPage   int  `json:"total_page"`
	CurrentPage int  `json:"current_page"`
	NextPage    int  `json:"next_page"`
	PrevPage    int  `json:"prev_page"`
	IsFirstPage bool `json:"is_first_page"`
	IsLastPage  bool `json:"is_last_page"`
}

// Set ... ページネーションを設定する
func (m *Object) Set(totalCount int) {
	if totalCount == 0 {
		return
	}
	m.TotalCount = totalCount
	m.OffsetCount = ((m.CurrentPage - 1) * m.PerPage)
	m.TotalPage = int(math.Ceil(float64(totalCount) / float64(m.PerPage)))
	if m.CurrentPage < m.TotalPage {
		m.NextPage = m.CurrentPage + 1
	}
	if 1 < m.CurrentPage {
		m.PrevPage = m.CurrentPage - 1
	}
	if m.TotalPage <= m.CurrentPage {
		m.IsLastPage = true
	}
	if m.CurrentPage == 1 {
		m.IsFirstPage = true
	}
	return
}

// New ... ページネーションを作成する
func New(currentPage int, perPage int) *Object {
	return &Object{
		CurrentPage: currentPage,
		PerPage:     perPage,
		IsFirstPage: false,
		IsLastPage:  false,
	}
}
