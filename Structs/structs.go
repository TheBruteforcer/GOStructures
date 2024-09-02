package structs

type Student struct {
	ID             int `gorm:"primaryKey;autoIncrement"`
	Name           string
	Code           int
	Rank           int
	AttendanceRate int
	Messages       []Messages
	Degrees        []Degrees
}

type Messages struct {
	ID        int `gorm:"primaryKey"`
	Content   string
	StudentID int    // Foreign key to Student
	Type      string // thank , dep

}

type Degrees struct {
	ID         int `gorm:"primaryKey"`
	TestTitle  string
	TestDegree int
	StudentID  int // Foreign key to Student
}

type Posts struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	Genre        string
	Title        string
	Description  string
	EmbededLinks string
}
