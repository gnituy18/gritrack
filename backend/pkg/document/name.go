package document

type Name string

const (
	Mission Name = "Mission"
	User         = "User"
)

func (n Name) String() string {
	return string(n)
}
