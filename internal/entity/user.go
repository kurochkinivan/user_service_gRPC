package entity

type User struct {
	ID        string
	Name      string
	Age       int32
	Gender    string
	About     string
	Photos    []*Photo
	Interests []*Interest
}
