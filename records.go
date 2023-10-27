package records

type Records interface {
	Connect() error
	Close() error

	CreateRecord(string, map[string]interface{}) (string, error)
	UpdateRecord(string, map[string]interface{}) (string, error)
	DeleteRecord(string) (string, error)

	FetchRecord(string) (string, map[string]interface{}, error)
}
