package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "Missing@2022"
	DB_NAME     = "postgres"
	host        = "localhost"
	port        = 5432
)

type Employee struct {
	EmployeeID   int    `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	ProjectName  string `json:"projectName"`
	Salary       int    `json:"salary"`
}

type JsonResponse struct {
	Type    string     `json:"type"`
	Data    []Employee `json:"data"`
	Message string     `json:"message"`
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	fmt.Println("DB info ::" + dbinfo)
	db, err := sql.Open("postgres", dbinfo)
	v, ok := db.Exec("SELECT * FROM employee.Employee")
	fmt.Println("pinging to DB", db.Ping())
	fmt.Println(ok)
	fmt.Println(v)
	checkErr(err)

	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/getemployees", getEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/employees", CreateEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", UpdateEmployees).Methods("PUT")
	router.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/employees/", DeleteEmployees).Methods("DELETE")
	fmt.Println("Server at 8080")
	http.ListenAndServe(":8000", router)
	//log.Fatal(http.ListenAndServe(":9000", router))

}
func main() {
	//setupDB()
	initializeRouter()

}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting employees...")

	// Get all employees from Employee table
	rows, err := db.Query("SELECT * FROM employee.Employee")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var employees []Employee

	// Foreach employee
	for rows.Next() {
		var employeeId int
		var employeeName string
		var projectName string
		var salary int

		err = rows.Scan(&employeeId, &employeeName, &projectName, &salary)

		// check errors
		checkErr(err)

		employees = append(employees, Employee{EmployeeID: employeeId, EmployeeName: employeeName, ProjectName: projectName, Salary: salary})
	}
	println(employees)
	var response = JsonResponse{Type: "success", Data: employees}

	json.NewEncoder(w).Encode(response)
}

// create employee

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee

	json.NewDecoder(r.Body).Decode(&employee)
	employeeID := employee.EmployeeID
	employeeName := employee.EmployeeName
	projectName := employee.ProjectName
	salary := employee.Salary

	if employeeID == 0 || employeeName == "" {
		response := JsonResponse{Type: "error", Message: "You are missing employeeId or employeename parameter."}
		json.NewEncoder(w).Encode(response)
	} else {
		db := setupDB()
		json.NewEncoder(w).Encode(&employee)

		printMessage("Inserting employee into DB")

		fmt.Println("Inserting new employee with ID: ", employeeID, " and name: ", employeeName)

		_, err := db.Query("INSERT INTO employee.Employee(employeeid, employeeName, projectName, salary) VALUES($1, $2, $3, $4);", employeeID, employeeName, projectName, salary)

		// check errors
		checkErr(err)

		response := JsonResponse{Type: "success", Message: "The employee record has been inserted successfully!"}
		json.NewEncoder(w).Encode(response)
	}
}

//deleting employee record

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	var employee Employee

	json.NewDecoder(r.Body).Decode(&employee)
	employeeID := employee.EmployeeID

	var response = JsonResponse{}

	if employeeID == 0 {
		response = JsonResponse{Type: "error", Message: "You are missing employee id parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting employee from DB")

		_, err := db.Exec("DELETE FROM employee.Employee where employeeid = $1", employeeID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The employee record has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteEmployees(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all employees...")

	_, err := db.Exec("DELETE FROM employee.Employee")

	// check errors
	checkErr(err)

	printMessage("All employees have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All employees have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}

func UpdateEmployees(w http.ResponseWriter, r *http.Request) {
	var employee Employee

	json.NewDecoder(r.Body).Decode(&employee)
	employeeID := employee.EmployeeID
	employeeName := employee.EmployeeName
	projectName := employee.ProjectName
	salary := employee.Salary

	if employeeID == 0 || employeeName == "" {
		response := JsonResponse{Type: "error", Message: "You are missing employeeId or employeename parameter."}
		json.NewEncoder(w).Encode(response)
	} else {
		db := setupDB()
		json.NewEncoder(w).Encode(&employee)

		printMessage("updating employee into DB")

		fmt.Println("updating new employee data with ID: ", employeeID, " and name: ", employeeName)

		_, delerr := db.Exec("DELETE FROM employee.Employee where employeeid = $1", employeeID)

		checkErr(delerr)

		_, err := db.Query("INSERT INTO employee.Employee(employeeid, employeeName, projectName, salary) VALUES($1, $2, $3, $4);", employeeID, employeeName, projectName, salary)

		// check errors
		checkErr(err)

		response := JsonResponse{Type: "success", Message: "The employee record has been inserted successfully!"}
		json.NewEncoder(w).Encode(response)
	}
}

func getEmployee(w http.ResponseWriter, r *http.Request) {

	var employee Employee

	json.NewDecoder(r.Body).Decode(&employee)
	employeeID := employee.EmployeeID

	if employeeID == 0 {
		response := JsonResponse{Type: "error", Message: "You are missing employeeId parameter."}
		json.NewEncoder(w).Encode(response)
	} else {
		db := setupDB()

		printMessage("get employee details from DB")

		_, err := db.Exec("SELECT * FROM employee.Employee where employeeid = $1", employeeID)

		// check errors
		checkErr(err)

		//response = JsonResponse{Type: "success", Message: "The employee record has been deleted successfully!"}
		json.NewEncoder(w).Encode(employee)
	}

}
