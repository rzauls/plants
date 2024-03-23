package httpd

import "fmt"

type routePathGenerator struct {
	root string
}

func (rpg *routePathGenerator) route(method, path string) string {
	return fmt.Sprintf("%s %s%s", method, rpg.root, path)
}
