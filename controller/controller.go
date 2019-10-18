// Package controller Test Project API's
//
// The purpose of this application is to store and retrieve test records for published test results
//
//
//
//     BasePath: /
//     Version: 1.0.0
//     License: bleh
//
//     Contact: Eshan Kaley<eshkaley@in.ibm.com>
//
//     Consumes:
//       - application/json
//
//     Produces:
//       - application/json
//
//     Security:
//       - token:
//
//     SecurityDefinitions:
//       token:
//         type: apiKey
//         in: header
//         name: Authorization
//
//
// swagger:meta
package controller

import (
	"fmt"
	"strconv"

	//"fmt"

	plane "RESTApp/dao/plane"
	student "RESTApp/dao/student"
	"RESTApp/model"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"

	//"errors"
	"encoding/json"
	"log"
)

//Handlers ...
func Handlers(ds *mgo.Session) http.Handler {
	server := mux.NewRouter() //create a new Server and attach handlers to it

	server.PathPrefix("/public/").Handler(
		http.StripPrefix("/public/", http.FileServer(http.Dir("/home/wiz/go/src/RESTApp/public/"))))

	server.HandleFunc("/", redir).Methods("GET")
	server.HandleFunc("/swagger", GetSwagger).Methods("GET")
	server.HandleFunc("/plane", AddPlane(ds)).Methods("POST")
	server.HandleFunc("/planes", GetPlanesHandler(ds)).Methods("GET")
	server.HandleFunc("/plane/{name}", RemovePlaneByName(ds)).Methods("DELETE")
	server.HandleFunc("/plane/{id}", RemovePlaneByID(ds)).Methods("DELETE")
	server.HandleFunc("/studentAggregates", StudentAggregates(ds)).Methods("GET")

	//Student Handlers
	server.HandleFunc("/student/{name}", DeleteStudent(ds)).Methods("DELETE") //done
	server.HandleFunc("/students", GetAllStudents(ds)).Methods("GET")         //done
	server.HandleFunc("/student/{name}", GetByName(ds)).Methods("GET")        //done
	server.HandleFunc("/student/{name}", UpdateStud(ds)).Methods("PUT")
	server.HandleFunc("/student", AddStudent(ds)).Methods("POST") // done        //done

	//Book Handlers
	//server.HandleFunc("/book", book.GetBookSession).Methods("GET")

	return server
}

// GetSwagger ...
func GetSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/swagger.json")
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/public/dist/#/", http.StatusFound)
}

//GetPlanesHandler ...
func GetPlanesHandler(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation GET /planes GET getPlanes
	//
	// Get Planes
	//
	// Get Catalog of planes
	// ---
	// produces:
	// - application/json
	// responses:
	//  '200':
	//    description: Found Results
	//    schema:
	//     type: array
	//     items:
	//      "$ref": "#/definitions/GetPlanesAPIResponse"
	//  '401':
	//    description: Unauthorized, Likely Invalid or Missing Token
	//  '403':
	//    description: Forbidden, you are not allowed to undertake this operation
	//  '404':
	//    description: Not found
	//  '500':
	//    description: Error occurred while processing the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			res, err := json.Marshal("Bad Request")
			if err != nil {
				log.Printf("Error while Marshalling: %v", err)
				sendErr(w, http.StatusMethodNotAllowed, res)
			}
		}

		allPlanes, err := plane.GetAllPlanes(ds)
		if err != nil {
			log.Printf("Error while Fetching Planes : %v ", err)
		}

		res, err := json.Marshal(allPlanes)
		if err != nil {
			log.Printf("Error while Marshalling to send Result")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})
}

// sendErr helper function
func sendErr(w http.ResponseWriter, stat int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	w.Write(res)
}

// AddPlane ...
func AddPlane(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation POST /plane POST putPlane
	//
	//
	// Put Plane in the Plane catalog
	//
	//
	// Add a New Plane
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: Plane
	//   type: object
	//   description: Student to be Added in the Catalog
	//   required: true
	//   in: body
	//   schema:
	//    $ref: '#/definitions/Plane'
	// responses:
	//  '200':
	//    description: Added Plane To the Catalog Successfully
	//    schema:
	//     $ref: '#/definitions/Plane'
	//  '401':
	//    description: Unauthorized, Likely Invalid or Missing Token
	//  '403':
	//    description: Forbidden, you are not allowed to undertake this operation
	//  '404':
	//    description: Not found
	//  '500':
	//    description: Error occurred while processing the request
	//
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			res, err := json.Marshal("Bad Request")
			if err != nil {
				log.Printf("Error while Marshalling : %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(res)
		}

		if r.Body != nil {
			var pl model.Plane
			err := json.NewDecoder(r.Body).Decode(&pl)
			if err != nil {
				log.Printf("Error while Decode body : %v", err)
			}
			plane.PutPlane(pl, ds)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			res, err := json.Marshal("Added Plane Successfully")
			if err != nil {
				log.Printf("Error while Marshalling : %v", err)
			}
			w.Write(res)
			defer r.Body.Close()
		}
	})
}

// RemovePlaneByName ...
func RemovePlaneByName(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation DELETE /plane/{name} DELETE removePlane
	//
	// Delete Plane
	//
	// Delete a Plane from Plane Catalog
	// ---
	// produces:
	// - application/json
	// - application/xml
	// parameters:
	// - name: name
	//   in: query
	//   required: true
	//   description: The name of the Plane to be removed
	// responses:
	//  '200':
	//    description: Plane Removed Successfully
	//  '401':
	//    description: Unauthorized, Likely Invalid or Missing Token
	//  '403':
	//    description: Forbidden, you are not allowed to undertake this operation
	//  '404':
	//    description: Not found
	//  '500':
	//    description: Error occurred while processing the request
	//
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			res, err := json.Marshal("Bad Request")
			if err != nil {
				log.Printf("Error while encoding error message : %v", err)
			}
			sendErr(w, http.StatusMethodNotAllowed, res)
		}
		params := mux.Vars(r)
		del := params["name"]
		ok := plane.DeletePlane(del, ds)
		if !ok {
			res, err := json.Marshal("Could Not Delete Server Error")
			if err != nil {
				log.Printf("Error while encoding error message : %v", err)
			}
			sendErr(w, http.StatusInternalServerError, res)
		}
		res, _ := json.Marshal("Deleted Successfully")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})
}

// RemovePlaneByID ...
func RemovePlaneByID(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation DELETE /plane/{id} DELETE removePlane
	//
	// Remove Plane
	//
	// Removes a Plane from DB
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   type: integer
	//   description: id of the plane to remove
	//   required: true
	//   in: query
	// responses:
	//   '200':
	//     description: Plane Removed Successfully
	//   '401':
	//     description: Unauthorized, Likely Invalid or Missing Token
	//   '403':
	//     description: Forbidden, you are not allowed to undertake this operation
	//   '404':
	//     description: Not found
	//   '500':
	//     description: Error occurred while processing the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			res, err := json.Marshal("Bad Request")
			if err != nil {
				log.Printf("Error while encoding error message : %v", err)
			}
			sendErr(w, http.StatusMethodNotAllowed, res)
		}

		countCheckSession := ds.Clone()
		count, err := countCheckSession.DB("trial").C("Student").Count()
		if err != nil {
			log.Println("Error while getting count from DB : %v", err)
		}

		params := mux.Vars(r)
		idIns := params["id"]
		id, err := strconv.Atoi(idIns)
		if err != nil {
			res, _ := json.Marshal("Server Error")
			sendErr(w, http.StatusInternalServerError, res)
		}
		if id > count {
			res, err := json.Marshal("ID with the given value Doesn't Exist")
			if err != nil {
				log.Printf("Error while encoding error message : %v", err)
			}
			sendErr(w, http.StatusBadRequest, res)
		}

		ok := plane.DeletePlaneByID(id, ds)
		if !ok {
			res, err := json.Marshal("Could Not Delete Server Error")
			if err != nil {
				log.Printf("Error while encoding error message : %v", err)
			}
			sendErr(w, http.StatusInternalServerError, res)
		}
		res, _ := json.Marshal("Deleted Successfully")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	})
}

//UpdateStud Update Student Info ByName
func UpdateStud(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation PUT /student/{name} UPDATE updateStudent
	//
	// Update a Students Information in The Student Catalog
	// ---
	// description: "Update Student Details"
	// summary: "Update Student Details in the Catalog"
	// parameters:
	// - name: name
	//   in: path
	//   description: "Student Name to Update Details"
	//   required: true
	//   type: string
	// - name : Student
	//   in: body
	//   required: true
	//   schema:
	//    $ref: '#/definitions/Student'
	// responses:
	//  '200':
	//   description: "Student Updated Successfully"
	//  '400':
	//   description: "Invalid Student Name Specified"
	//  '404':
	//   description: "Student Not Found"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check for method PUT, if anything else, respond with appropriate status
		if r.Method != "PUT" {
			res, err := json.Marshal("Bad Request")
			if err != nil {
				log.Printf("Bad Request : %v", err)
			}

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(res)
		}

		if r.Body != nil {
			var stuNew model.Student //student object to store updated student

			//extract name from path
			params := mux.Vars(r)
			nm := params["name"]
			fmt.Println(nm)
			defer r.Body.Close()

			//get a studentObject from GetByName using the extracted name
			stuToChange, err := student.GetByName(nm, ds)
			if err != nil {
				log.Printf("Error While Fetching Record to Update : %v", err)
				res, _ := json.Marshal("Invalid Student Name specified")
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write(res) //send Message
				return

			}

			//Decode values from body sent from client into a studentObject
			err = json.NewDecoder(r.Body).Decode(&stuNew)
			if err != nil {
				log.Printf("Error While Deconding Body : %v", err)
			}

			//update the values from the body to The object got from the Database
			//stuToChange.StudentName = stuNew.StudentName
			stuToChange.StudentAge = stuNew.StudentAge     //update age
			stuToChange.StudentMarks = stuNew.StudentMarks //update marks

			//respond with appropriate message after calling Data Access Layer
			err = student.UpdateStudent(stuToChange, ds)
			if err != nil {
				log.Printf("Could not Update student: %v", err)
			}

			res, err := json.Marshal("Updated Successfully")
			if err != nil {
				log.Printf("Error While Marshalling! : %v", err)
			}

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res) //send Message
		}

	})
}

//GetByName ...
func GetByName(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation GET /student/{name} GET getStudent
	//
	// ---
	// description: Get Student Details by name
	// summary: "Get Student By Name"
	// parameters:
	// - in: path
	//   name: name
	//   required: true
	//   schema:
	//    type: string
	// responses:
	//  '200':
	//   description: Details Fetched Status Ok
	//   schema:
	//    type: object
	//    $ref: '#/definitions/Student'
	//    example:
	//     studentName: Eshan
	//     studentAge: 25
	//     studentMarks: 70
	//  '400':
	//     description: "No Entry Found By that Name"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check for method GET, if anything else, respond with appropriate status
		if r.Method != "GET" {

			res, err := json.Marshal("Bad Request")

			if err != nil {
				log.Fatal(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, err.Error(), 200)
			w.Write(res)
		}

		params := mux.Vars(r) //extract name from URL path

		var s model.Student
		s, err := student.GetByName(params["nm"], ds) //call data access layer

		if err != nil {
			res, _ := json.Marshal("No Entry Found By That Name")
			w.Header().Set("Content Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(res)
			return
		}

		//respond with appropriate message
		mresult, _ := json.Marshal(s)
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mresult)
	})
}

//AddStudent ...
func AddStudent(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation POST /student POST AddStudent
	//
	// Add Student
	//
	// Add a Student to the Student
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: Student
	//   required: true
	//   in: body
	//   schema:
	//    "$ref": '#/definitions/Student'
	// responses:
	//  '200':
	//   description: Added Student Successfully to the Catalog
	//  '401':
	//   description: Unauthorized, Likely Invalid or Missing Token
	//  '403':
	//   description: Forbidden, you are not allowed to undertake this operation
	//  '404':
	//   description: Not found
	//  '500':
	//   description: Error occurred while processing the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check if method is POST else show error
		if r.Method != "POST" {

			response, _ := json.Marshal("Bad Request")
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(response)
		}
		//check if body has content
		if r.Body != nil {
			defer r.Body.Close()
			var stu model.Student

			//decode the body for student details
			err := json.NewDecoder(r.Body).Decode(&stu)
			if err == nil {
				student.AddStudent(stu, ds)
				//w.Header().Set("Access-Control-Allow-Methods","POST,OPTIONS")
				response, _ := json.Marshal("Added Successfully")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			}
		}
	})
}

//DeleteStudent ...
func DeleteStudent(ds *mgo.Session) http.HandlerFunc {
	// swagger:operation DELETE /student/{name} DELETE deleteStudent
	//
	// Delete Stident
	//
	// Delete A Student From Student Catalog
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: name
	//   type: string
	//   description: Name of the Student to Delete
	//   in: query
	//   required: true
	// responses:
	//  200:
	//   description: Removed Student from the Catalog
	//  400:
	//   description: "Invalid Student Name Specified"
	//  404:
	//   description: "Student Not Found"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if method is DELETE else respond with error
		if r.Method != "DELETE" {
			response, err := json.Marshal("Bad Request")
			if err != nil {
				log.Printf("Bad Request!: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(response)

		}

		//check if body has content
		if r.Body != nil {

			defer r.Body.Close()
			params := mux.Vars(r) //extract name of student from URL path

			//err := dao.GetByName(params["name"])
			//Respond to the requeset after calling Data Access Layer
			err := student.RemoveByName(params["name"], ds)
			if err != nil {
				res, _ := json.Marshal("Could not Find anyone with that name")
				w.Header().Set("Content-Type", "appication/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write(res)
				return
			}
			response, err := json.Marshal("Removed Student")
			if err != nil {
				log.Fatal(err)
				return
			}

			w.Header().Set("Content-Type", "appication/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	})
}

//GetAllStudents ...
func GetAllStudents(ds *mgo.Session) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// swagger:operation GET /students GET getAllStudents
		//
		// List of Students
		//
		// Get the Student Catalog in Response
		//
		// ---
		// produces:
		// - application/json
		// responses:
		//  '200':
		//   description: Found Results
		//   schema:
		//    type: array
		//    items:
		//     "$ref": "#/definitions/GetAllStudentsAPIResponse"
		//  '400':
		//   description: "Invalid Student Name Specified"
		//  '404':
		//   description: "Student Not Found"

		//check for method GET, if any other, respond with error with appropriate status
		if r.Method != "GET" {
			response, err := json.Marshal("Bad Request")

			if err != nil {
				log.Printf("Bad Request!: %v ", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			//w.Header().S
			//w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(response)
		}

		//respond with appropriate message after calling Data Access Layer
		res, err := student.GetAll(ds)
		if err != nil {
			log.Fatal(err)
		}
		response, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	})
}
