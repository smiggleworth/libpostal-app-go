package main

import (
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
)

func toCamelInitCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := false
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

func parse(a string) map[string]string {
	m := make(map[string]string)
	for _, c := range parser.ParseAddress(a) {
		m[toCamelInitCase(c.Label)] = c.Value
	}
	return m
}

func explode(a string) []map[string]string {
	r := make([]map[string]string, 0)
	for _, entry := range expand.ExpandAddress(a) {
		r = append(r, parse(entry))
	}
	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"api": "libpostal",
			"ver": "v1.0.0"})
	})

	r.GET("/parse", func(c *gin.Context) {
		a := c.Query("address")
		if a == "" {
			c.JSON(400, gin.H{"message": "invalid address"})
			return
		}
		c.JSON(200, parse(a))
	})

	r.GET("/expand", func(c *gin.Context) {
		a := c.Query("address")
		if a == "" {
			c.JSON(400, gin.H{"message": "invalid address"})
			return
		}
		c.JSON(200, expand.ExpandAddress(a))
	})

	r.GET("/explode", func(c *gin.Context) {
		a := c.Query("address")
		if a == "" {
			c.JSON(400, gin.H{"message": "invalid address"})
			return
		}
		c.JSON(200, explode(a))
	})

	r.Run(":80")
}
