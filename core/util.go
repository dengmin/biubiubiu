package core

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"net/http"
)

func md5Encrypt(data string) string{
	h := md5.New()
	h.Write([]byte(string(data)))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func getRealIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func constantIp(ips []string, realIp string) bool {
	for i := 0; i < len(ips); i++ {
		if ips[i] == realIp {
			return true
		}
	}
	return false
}
