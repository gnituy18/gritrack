package document

type Name string

const (
	Mission Name = "Mission"
	User         = "User"
	Step         = "Step"
)

func (n Name) String() string {
	return string(n)
}
