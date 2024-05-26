package middleware

import (
	"net"
	"net/http"
)

type Checker struct {
	subnet net.IPNet
}

func NewChecker(subnet net.IPNet) *Checker {
	return &Checker{subnet: subnet}
}

func CheckSubnet(subnet net.IPNet) func(next http.Handler) http.Handler {
	c := NewChecker(subnet)
	return c.CheckHandler
}

func (c *Checker) CheckHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			http.Error(w, "No real ip provided", http.StatusForbidden)
			return
		}
		ip := net.ParseIP(realIP)
		if ip == nil {
			http.Error(w, "Bad X-Real-IP header", http.StatusForbidden)
			return
		}
		if !c.subnet.Contains(ip) {
			http.Error(w, "Ip is not in trusted network", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
