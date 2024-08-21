package entity

type Config struct {
    MaxCustomer int `json:"max_customer" validate:"required,min=2"`
}