package utils

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"treehole_next/models"
)

//func InArray[T comparable](item *T, container *[]T) bool {
//	for _, i := range *container {
//		if *item == i {
//			return true
//		}
//	}
//	return false
//}

var emptyMap = models.Map{}

// BindJSON is a safe method to bind request body to struct
func BindJSON(c *fiber.Ctx, obj interface{}) error {
	body := c.Body()
	if len(body) == 0 {
		body, _ = json.Marshal(emptyMap)
	}
	return json.Unmarshal(body, obj)
}
