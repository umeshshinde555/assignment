package main

import (
	"fmt"
	"net/http"
	"src/assignment/models"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	employees     = make(map[int]*models.Employee)
	lock          = sync.RWMutex{}
	lastID    int = 0
)

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Printf("Received %s request for %s\n", ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Next()
	}
}

func getRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/employees", listEmployees)
	router.POST("/employees", createEmployee)
	router.GET("/employees/:id", getEmployee)
	router.PUT("/employees/:id", updateEmployee)
	router.DELETE("/employees/:id", deleteEmployee)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": 404, "error": "Not Found"})
	})
	return router
}
func main() {
	router := getRoutes()
	router.Run(":8080")
}

func listEmployees(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	start := (page - 1) * pageSize
	end := start + pageSize

	var slice []*models.Employee
	i := 0
	for _, emp := range employees {
		if i >= start && i < end {
			slice = append(slice, emp)
		}
		i++
	}

	c.JSON(http.StatusOK, slice)
}

func createEmployee(c *gin.Context) {
	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := CreateEmployeeFun(&emp)
	emp.ID = id
	c.JSON(http.StatusCreated, emp)
}

func getEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	if emp, exists := GetEmployeeByID(id); exists {
		c.JSON(http.StatusOK, emp)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
	}
}

func updateEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if UpdateEmployeeFun(id, &emp) {
		c.JSON(http.StatusOK, emp)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
	}
}

func deleteEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	if DeleteEmployeeFun(id) {
		c.JSON(http.StatusOK, gin.H{"message": "employee deleted"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
	}
}

func CreateEmployeeFun(e *models.Employee) int {
	lock.Lock()
	defer lock.Unlock()
	lastID++
	e.ID = lastID
	employees[e.ID] = e
	return e.ID
}

func GetEmployeeByID(id int) (*models.Employee, bool) {
	lock.RLock()
	defer lock.RUnlock()
	emp, exists := employees[id]
	return emp, exists
}

func UpdateEmployeeFun(id int, e *models.Employee) bool {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := employees[id]; exists {
		e.ID = id
		employees[id] = e
		return true
	}
	return false
}

func DeleteEmployeeFun(id int) bool {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := employees[id]; exists {
		delete(employees, id)
		return true
	}
	return false
}
