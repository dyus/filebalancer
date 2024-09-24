package balancer

type Balancer interface {
	GetHosts(int) []string
	AddHost(string)
}
