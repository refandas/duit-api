package domain

type User struct {
	Id        string `dynamodbav:"Id"`
	Name      string `dynamodbav:"Name"`
	Email     string `dynamodbav:"Email"`
	Password  string `dynamodbav:"Password"`
	CreatedAt int64  `dynamodbav:"CreatedAt"`
}
