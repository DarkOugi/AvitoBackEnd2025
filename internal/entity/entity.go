package entity

type Merch struct {
	Name string
	Cnt  int
}

type User struct {
	Name     string
	Cost     int
	Password string
	// transaction float64 // в задании не используется, но в целом позитивное поле
}
