package amqp

type TCPSection struct {
	host string
	port int32
	user string
	pwd  string
	path string
}

func (m *TCPSection) GetHost() string {
	return m.host
}

func (m *TCPSection) GetPort() int32 {
	return m.port
}

func (m *TCPSection) GetUser() string {
	return m.user
}

func (m *TCPSection) GetPwd() string {
	return m.pwd
}

func (m *TCPSection) GetPath() string {
	return m.path
}
