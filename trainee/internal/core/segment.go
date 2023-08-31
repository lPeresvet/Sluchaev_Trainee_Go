package core

type Segment struct {
	Slug string `json:"slug"`
}

type UserSegmentId struct {
	UserId    int64
	SegmentId int64
}
