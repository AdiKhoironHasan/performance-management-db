package web

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"testing"
)

type Person struct {
	ID    int
	Name  string
	Email string
}

func TestServe(t *testing.T) {
	people := []Person{
		{1, "John Doe", "john@example.com"},
		{2, "Jane Smith", "jane@example.com"},
		{3, "Bob Johnson", "bob@example.com"},
	}

	// Create dynamic column names
	var columns []string
	if len(people) > 0 {
		columns = getColumnNames(people[0])
	}

	// Define the template
	tmpl := template.Must(template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Dynamic Table</title>
		</head>
		<body>
			<h1>Dynamic Table</h1>
			<table>
				<tr>
					{{range .Columns}}
						<th>{{.}}</th>
					{{end}}
				</tr>
				{{range .People}}
					<tr>
						{{range $index, $value := .}}
							<td>1</td>
						{{end}}
					</tr>
				{{end}}
			</table>
		</body>
		</html>
	`))

	// Define the handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			People  []Person
			Columns []string
		}{
			People:  people,
			Columns: columns,
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	})

	// Start the server
	log.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Helper function to get column names from a struct
func getColumnNames(data interface{}) []string {
	var columns []string
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		columns = append(columns, t.Field(i).Name)
	}
	return columns
}
