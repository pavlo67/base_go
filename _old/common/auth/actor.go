package auth

type Actor struct {
	*Identity `json:"identity,omitempty" bson:"identity,omitempty"`
	Creds     `json:"creds,omitempty"    bson:"creds,omitempty"`
}
