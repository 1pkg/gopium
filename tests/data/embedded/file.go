//+build tests_data

package embedded

import "time"

type MetaLabaratory struct {
}

type Person struct {
	Birtday time.Time `json:"birthday" db:"birthday"`
	Weight  float64   `json:"weight" db:"weight"`
	Height  float64   `json:"height" db:"height"`
}

type PatientObject struct {
	MetaLabaratory
	Person
	ID           string  `json:"id" db:"id"`
	Enrolled     bool    `json:"enrolled" db:"enrolled"`
	Gender       string  `json:"gender" db:"gender"`
	PhoneNumber  *string `json:"phone_number" db:"phone_number"`
	Email        *string `json:"email" db:"email"`
	AddressTitle *string `json:"address_title" db:"address_title"`
}
