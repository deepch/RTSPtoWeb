package main

import "github.com/gin-gonic/gin"

func HTTPAPIServerStreams(c *gin.Context) {
	c.IndentedJSON(200, Storage.List())
}

func HTTPAPIServerStreamAdd(c *gin.Context) {
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	err = Storage.StreamAdd(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}
func HTTPAPIServerStreamEdit(c *gin.Context) {
	if !Storage.StreamExist(c.Param("uuid")) {
		return
	}
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	err = Storage.StreamEdit(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}
func HTTPAPIServerStreamDelete(c *gin.Context) {
	err := Storage.StreamDelete(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}
func HTTPAPIServerStreamReload(c *gin.Context) {
	err := Storage.StreamReload(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}

func HTTPAPIServerStreamInfo(c *gin.Context) {
	info, err := Storage.StreamInfo(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, info)
}
