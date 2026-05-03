package entities

type MediaSource struct {
	ID         int
	FullName   string
	URL        string
	TrustLevel string
	IsBlocked  bool
}

type Publication struct {
	ID       int
	Title    string
	URL      string
	SourceID int
	Status   string
}

type Incident struct {
	ID            int
	Category      string
	Description   string
	Severity      string
	PublicationID int
}
