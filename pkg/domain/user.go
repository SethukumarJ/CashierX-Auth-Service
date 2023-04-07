package domain

type Users struct {
	Id        int64  `json:"id" gorm:"primaryKey"`
	UserName string `json:"user_name" gorm:"unique" validate:"required,min=2,max=50"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"notnull;unique" validate:"email,required"`
	Password  string `json:"password"`
}
