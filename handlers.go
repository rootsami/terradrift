package main

import "github.com/gin-gonic/gin"

func scanHandler(c *gin.Context) {

	name := c.Query("stack")
	planResp, err := stackScan(name)

	if err == nil {

		c.JSON(200, planResp)
	} else {
		// TODO: it looks ugly but cannot compare err with a string. has to be converted to string then pass it to handler. there has to be a better way
		errorMessage := error.Error(err)
		if errorMessage == "ERROR: STACK WAS NOT FOUND" {
			c.JSON(404, errorMessage)
		} else {
			c.JSON(500, errorMessage)
		}
	}
}
