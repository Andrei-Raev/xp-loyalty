package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Msg struct {
	Msg string `json:"message"`
}

type Err struct {
	Err string `json:"error"`
}

func E(err error) *Err {
	return &Err{
		Err: err.Error(),
	}
}

func M(msg string) *Msg {
	return &Msg{
		Msg: msg,
	}
}

func ParsePath(c *gin.Context, param string) (string, error) {
	p := c.Param(param)
	if p == "" {
		return "", errors.New("empty param")
	}
	return p, nil
}
