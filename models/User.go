package models


import(
	"time"
	"database/sql"
)

type User struct {	
	Id  		int32   		  `json:"id"` //,omitempty
	Username 	string    		  `json:"username"` //,omitempty
	Password	string    		  `json:"password"` //,omitempty
	Authority 	int32    		  `json:"authority"` //,omitempty
	Created_on	time.Time    	  `json:"created_on"` //,omitempty
	Last_login	sql.NullTime	   `json:"last_login"` //,omitempty
}
