package models

// PermType2

import "encoding/json"

// PermType2
type PermType2 struct {
	Owner        string       `json:"owner"`
	OwnerAccess  AccessType   `json:"owner_access"`
	GlobalAccess AccessType   `json:"global_access"`
	Share        []*ShareType `json:"share"`
}

// String returns json representation of the object
func (model *PermType2) String() string {
	b, _ := json.Marshal(model)
	return string(b)
}

// MakePermType2 makes PermType2
func MakePermType2() *PermType2 {
	return &PermType2{
		//TODO(nati): Apply default
		Owner:        "",
		OwnerAccess:  MakeAccessType(),
		GlobalAccess: MakeAccessType(),

		Share: MakeShareTypeSlice(),
	}
}

// MakePermType2Slice() makes a slice of PermType2
func MakePermType2Slice() []*PermType2 {
	return []*PermType2{}
}
