package main

// I hope you will consider my assessment even if there are any errors,
// I learnt GoLang to finish this assessment, I am very much interested in being a part of Fetch.
import (
	"Reciept_APIs/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Starting the engine.
// Defining the routes and endpoints.
func main() {
	server := gin.Default()
	server.GET("/reciept/:id/points", getPoints)
	server.POST("/reciepts/process", processReciept)
	server.Run(":8080")
}

// Endpoint which accepts the recipet JSON, calls the "ProcessReciept" and returns a JSON object of ID/key of it.
func processReciept(context *gin.Context) {
	var reciept models.Reciept
	err := context.ShouldBindBodyWithJSON(&reciept)
	// If condition to check if there is any content in the JSON object and it abides the Reciept struct.
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request."})
		return
	}
	// Calling the "ProcessReciept" service which returns an id.
	id := reciept.SaveReciept()
	// If condition to check if there is any error in processing the reciept, if any return the status and halt the process.
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Reciept could not be processed."})
		return
	}
	// If the reciept has been processed(saved) return JSON object of the id
	context.JSON(http.StatusCreated, gin.H{"Reciept ID": id})
}

// Endpoint which accepts the id(through the context) and calls the service "CalculatePoints" and returns the total points.
func getPoints(context *gin.Context) {
	// Varibale to store the retrieved parameter, "id".
	recieptId := context.Param("id")
	// Variable to store the boolean value if there is an id.
	_, isExists := models.AllReceipts[recieptId]
	// If condition to check if there isn't any id present, if true then returns a JSON object of http status bad request.
	if !isExists {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Id passed"})
		return
	}
	// If 'id' exists then, the variable 'points' would store and be returned as a JSON object.
	points := models.CalculatePoints(recieptId)
	context.JSON(http.StatusOK, gin.H{"total points": points})

}
