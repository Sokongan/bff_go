package domain

type ConsentRequest struct {
	Skip           bool
	ClientID       string
	RequestedScope []string
	Audience       []string
	Subject        string
}
