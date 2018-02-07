package geanlib

// AuthorList is a list of all authors and their metadata.
type AuthorList map[string]Author

// Author contains details about the author of a page.
type Author struct {
	GivenName   string
	FamilyName  string
	DisplayName string
	Thumbnail   string
	Image       string
	ShortBio    string
	LongBio     string
	Email       string
	Social      AuthorSocial
}

// AuthorSocial is a place to put social details per author. These are the
// standard keys that themes will expect to have available, but can be
// expanded to any others on a per site basis
// - website
// - github
// - facebook
// - twitter
// - googleplus
// - pinterest
// - instagram
// - youtube
// - linkedin
// - skype
type AuthorSocial map[string]string
