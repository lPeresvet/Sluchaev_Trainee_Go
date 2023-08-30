package core

type User struct {
	Id int64 `json:"id"`
	Segments []*Segment `json:"activeSegment,omitempty"`
}
