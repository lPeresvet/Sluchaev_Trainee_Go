package core

type UserRequest struct {
	Id int64 `json:"id"`
	SegmentsToAdd []*Segment `json:"segmentsToAdd,omitempty"`
	SegmentsToDelete []*Segment `json:"segmentsToDelete,omitempty"`
}