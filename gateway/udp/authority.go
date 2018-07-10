package udp

var authorities = map[string]Authority{
	"x": Authority{10000},
	"y": Authority{20000},
}

type Authority []int

func (auth Authority) Check(port int) bool {
	for i := range auth {
		if port == auth[i] {
			return true
		}
	}
	return false
}
