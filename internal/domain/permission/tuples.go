package permission_domain

type RelationTuple struct {
	Namespace string `json:"namespace"`
	Object    string `json:"object"`
	Relation  string `json:"relation"`
	SubjectID string `json:"subject_id"`
}

type ListTuplesParams struct {
	Namespace string
	Object    string
	Relation  string
	SubjectID string
	PageToken string
	PageSize  int64
}

type ListTuplesResult struct {
	Tuples        []RelationTuple `json:"tuples"`
	NextPageToken string          `json:"next_page_token,omitempty"`
}
