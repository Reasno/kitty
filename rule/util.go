package rule

import (
	"net/http"
	"strings"
)

//TrimPrefix removes the longest common Prefix from all provided strings
func TrimPrefix(strs []string) {
	p := Prefix(strs)
	if p == "" {
		return
	}
	for i, s := range strs {
		strs[i] = strings.TrimPrefix(s, p)
	}
}

//TrimSuffix removes the longest common Suffix from all provided strings
func TrimSuffix(strs []string) {
	p := Suffix(strs)
	if p == "" {
		return
	}
	for i, s := range strs {
		strs[i] = strings.TrimSuffix(s, p)
	}
}

//Prefix returns the longest common Prefix of the provided strings
func Prefix(strs []string) string {
	return longestCommonXfix(strs, true)
}

//Suffix returns the longest common Suffix of the provided strings
func Suffix(strs []string) string {
	return longestCommonXfix(strs, false)
}

func longestCommonXfix(strs []string, pre bool) string {
	//short-circuit empty list
	if len(strs) == 0 {
		return ""
	}
	xfix := strs[0]
	//short-circuit single-element list
	if len(strs) == 1 {
		return xfix
	}
	//compare first to rest
	for _, str := range strs[1:] {
		xfixl := len(xfix)
		strl := len(str)
		//short-circuit empty strings
		if xfixl == 0 || strl == 0 {
			return ""
		}
		//maximum possible length
		maxl := xfixl
		if strl < maxl {
			maxl = strl
		}
		//compare letters
		if pre {
			//Prefix, iterate left to right
			for i := 0; i < maxl; i++ {
				if xfix[i] != str[i] {
					xfix = xfix[:i]
					break
				}
				if len(xfix) > maxl {
					xfix = xfix[:maxl]
				}
			}
		} else {
			//Suffix, iternate right to left
			for i := 0; i < maxl; i++ {
				xi := xfixl - i - 1
				si := strl - i - 1
				if xfix[xi] != str[si] {
					xfix = xfix[xi+1:]
					break
				}
				if len(xfix) > maxl {
					xfix = xfix[len(xfix)-maxl:]
				}
			}
		}
	}
	return xfix
}

func DbKeys(m map[string]Container) []string {
	var out []string
	for _, v := range m {
		out = append(out, v.DbKey)
	}
	return out
}

func realIP(r *http.Request) string {
	var ip string
	var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
