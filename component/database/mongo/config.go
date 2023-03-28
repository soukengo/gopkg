package mongo

type Config struct {
	Address    string
	Username   string
	Password   string
	Database   string
	AuthSource string
	Timeout    int
}
