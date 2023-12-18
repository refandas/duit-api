package domain

// Spending represent the spending history data structure.
type Spending struct {

	// Id represents the unique identifier of a user's spending history.
	// It is formatted as a UUID4.
	Id string `dynamodbav:"Id"`

	// UserId represents the unique identifier of the user who made
	// the spending. It is formatted as a UUID4
	UserId string `dynamodbav:"UserId"`

	// Title represents the title or name associated with user's spending.
	Title string `dynamodbav:"Title"`

	// Description represents additional details or notes regarding
	// a user's spending.
	Description string `dynamodbav:"Description"`

	// Amount represents the spending amount. It is used to store the
	// monetary value of a user's spending.
	Amount float64 `dynamodbav:"Amount"`

	// Date represents the date when a spending was done, stored in
	// Unix time format. It is used to store the timestamp of when
	// the spending occurred.
	Date int64 `dynamodbav:"Date"`

	// Category represents the spending category chosen by the user.
	Category string `dynamodbav:"Category"`

	// CreatedAt represents the date and time when the spending data
	// was created, stored in Unix time format. It is used to store
	// the timestamp of when the spending data was initially recorded.
	CreatedAt int64 `dynamodbav:"CreatedAt"`
}
