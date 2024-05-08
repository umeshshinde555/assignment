package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"src/assignment/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/employees", listEmployees)
	r.POST("/employees", createEmployee)
	r.GET("/employees/:id", getEmployee)
	r.PUT("/employees/:id", updateEmployee)
	r.DELETE("/employees/:id", deleteEmployee)
	return r
}

func TestCreateEmployeeHandler(t *testing.T) {
	router := setupRouter()

	newEmployee := models.Employee{Name: "John Doe", Position: "Developer", Salary: 80000}
	data, _ := json.Marshal(newEmployee)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(data))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	var emp models.Employee
	json.Unmarshal(resp.Body.Bytes(), &emp)
	assert.Equal(t, "John Doe", emp.Name)
	assert.Equal(t, "Developer", emp.Position)
	assert.Equal(t, 80000.0, emp.Salary)
}

func TestGetEmployeeHandler(t *testing.T) {
	router := setupRouter()

	// Assuming ID 1 is present
	req, _ := http.NewRequest("GET", "/employees/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestListEmployeesHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/employees?page=1&pageSize=5", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUpdateEmployeeHandler(t *testing.T) {
	router := setupRouter()

	updatedEmployee := models.Employee{Name: "Jane Doe", Position: "Senior Developer", Salary: 100000}
	data, _ := json.Marshal(updatedEmployee)
	req, _ := http.NewRequest("PUT", "/employees/1", bytes.NewBuffer(data))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteEmployeeHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/employees/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
