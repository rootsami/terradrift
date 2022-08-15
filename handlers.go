package main

import "github.com/gin-gonic/gin"

func scanHandler(c *gin.Context) {

	name := c.Params.ByName("name")
	planResp, err := stackScan(name)

	if err == nil {

		c.JSON(200, planResp)
	} else {
		// IDK what's going here but finally it worked returning the error message
		errorMessage := error.Error(err)
		c.JSON(500, errorMessage)
	}
}
