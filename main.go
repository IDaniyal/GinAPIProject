package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID        string
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
	Age       int `json: "age"`
}

var UserList []User

func main() {
	router := gin.Default()

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", GetUsers)
		userRoutes.POST("/", AddUser)
		userRoutes.PUT("/:id", UpdateUser)
		userRoutes.DELETE("/:id", DeleteUser)
	}

	router.Run(":5050")
}

func GetUsers(c *gin.Context) {
	c.JSON(200, UserList)
}

func AddUser(c *gin.Context) {
	var reqBody User

	if error := c.ShouldBindJSON(&reqBody); error != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "invalid request body",
		})
		return
	}

	reqBody.ID = uuid.New().String()

	if !validation(&reqBody, c) {
		return
	}

	UserList = append(UserList, reqBody)
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Added User successfully",
	})
	return

}

func UpdateUser(c *gin.Context) {
	var reqBody User

	id := c.Param("id")

	if error := c.ShouldBindJSON(&reqBody); error != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "invalid request body",
		})
		return
	}

	for index, value := range UserList {
		if value.ID == id {
			if !validation(&reqBody, c) {
				return
			}
			UserList[index].FirstName = reqBody.FirstName
			UserList[index].LastName = reqBody.LastName
			UserList[index].Age = reqBody.Age
			c.JSON(200, UserList[index])
			return
		}
	}

	c.JSON(404, gin.H{
		"error":   true,
		"message": "User not found",
	})
	return
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	for index, value := range UserList {
		if value.ID == id {
			UserList = append(UserList[0:index], UserList[index+1:]...)
			c.JSON(200, gin.H{
				"error":   false,
				"message": "Deleted Successfully",
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"error":   true,
		"message": "User not found",
	})
	return

}


// Temporarily using a function for validation at the time of insertion
// and updating.

func validation(u *User, c *gin.Context) bool {
	if len(u.FirstName) == 0 || len(u.LastName) == 0 {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Length of name can not be 0",
		})
		return false
	}

	if u.Age == 0 {
		u.Age = 18
	}
	return true
}
