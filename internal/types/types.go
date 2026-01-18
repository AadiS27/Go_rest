package types 


type Student struct {
	Id int 
	Name string `validate:"required,min=2,max=100"`
	Email string `validate:"required,email"`
	Age int `validate:"required,min=18,max=100"`
}