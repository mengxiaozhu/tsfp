package tsfp

type Server struct {
	Addr    string
	Pattern string
	Init    string
}
type Config struct {
	Port    int
	Init    string
	Addr    string
	Static  string
	Servers map[string]Server
}

