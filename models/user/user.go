package user

type User struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Gender    string  `json:"gender"`
	Age       int     `json:"age"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	ID        int     `json:"id"`
	Email     string  `json:"email"`
	Password  string  `json:"password_hash"`
	Name      string  `json:"name"`
	Gender    string  `json:"gender"`
	Age       int     `json:"age"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	YesSwipes *int    `json:"yes_swipes"` // This field is only used when sorting by rank
}

type DiscoverUserResponseBody struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	Age            int    `json:"age"`
	DistanceFromMe int    `json:"distanceFromMe"`
}

type MatchedResultResponseBody struct {
	Match   bool `json:"match"`
	MatchID *int `json:"matchId"` // This field is only used when a match is created
}
