package main

type User struct {
	Name     string `json:"name" rh_label:"Name" rh_valid:"required,alpha"`
	UserCode string `json:"user_code" rh_label:"User Code" rh_valid:"required,len=5"`
	Username   string `json:"username" rh_label:"Username" rh_valid:"required,min=5,max=10"`
	Email      string `json:"email" rh_label:"Email" rh_valid:"required,email"`
}
