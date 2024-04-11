package user

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

type UserCreateResponseBody struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

type UserLoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponseBody struct {
	Token string `json:"token"`
}

type UsersTableRecord struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password_hash"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DiscoverUserResponseBody struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

type MatchedResultResponseBody struct {
	Match   bool `json:"match"`
	MatchID int  `json:"matchId"`
}

// continue the code here
