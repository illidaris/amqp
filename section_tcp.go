package amqp

type TCPSection struct {
	host string
	port int32
	user string
	pwd  string
	path string
}

func (m *TCPSection) SetHost(v string) {
	m.host = v
}

func (m *TCPSection) GetHost() string {
	return m.host
}

func (m *TCPSection) SetPort(v int32) {
	m.port = v
}

func (m *TCPSection) GetPort() int32 {
	return m.port
}

func (m *TCPSection) SetUser(v string) {
	m.user = v
}

func (m *TCPSection) GetUser() string {
	return m.user
}

func (m *TCPSection) SetPwd(v string) {
	m.pwd = v
}

func (m *TCPSection) GetPwd() string {
	return m.pwd
}

func (m *TCPSection) SetPath(v string) {
	m.path = v
}

func (m *TCPSection) GetPath() string {
	return m.path
}
