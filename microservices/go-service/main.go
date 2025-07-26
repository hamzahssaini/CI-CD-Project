package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var mongoStatus string

func connectMongo() {
	mongoURI := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		mongoStatus = "❌ MongoDB connection failed"
		log.Println(err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		mongoStatus = "❌ MongoDB ping failed"
		log.Println(err)
		return
	}

	collection = client.Database("myappdb").Collection("users")
	mongoStatus = "✅ MongoDB connected"
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	userCount := 0
	if collection != nil {
		count, err := collection.CountDocuments(context.TODO(), bson.D{})
		if err == nil {
			userCount = int(count)
		}
	}

	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Go Service – Hamza</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background: linear-gradient(135deg, #667eea, #764ba2);
                color: white;
                text-align: center;
                padding: 50px;
            }
            input, button {
                padding: 10px;
                border-radius: 5px;
                border: none;
                margin: 5px;
            }
            button {
                background: white;
                color: #764ba2;
                font-weight: bold;
                cursor: pointer;
            }
            a {
                color: #ffffff;
                font-weight: bold;
                text-decoration: underline;
            }
        </style>
    </head>
    <body>
        <h1>👋 Welcome to Go Service at <code>go.hamza.local</code></h1>
        <p>{{.MongoStatus}}</p>
        <p>📄 Total Subscribers: {{.UserCount}}</p>

        <form action="/register" method="POST">
            <input name="name" placeholder="Enter Your Name" required />
            <input name="email" type="email" placeholder="Enter Your Email" required />
            <button type="submit">Subscribe</button>
        </form>

        <br>
        <a href="/users">🔍 View all users</a>
    </body>
    </html>
    `
	tmplParsed := template.Must(template.New("home").Parse(tmpl))
	tmplParsed.Execute(w, map[string]interface{}{
		"MongoStatus": mongoStatus,
		"UserCount":   userCount,
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	_, err := collection.InsertOne(context.TODO(), bson.M{
		"name":   name,
		"email":  email,
		"source": "go-service",
	})

	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success?name="+name, http.StatusFound)
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, `<h2>✅ Thank you, <b>%s</b>! You're subscribed!</h2><a href="/">⬅️ Back</a>`, name)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if collection == nil {
		http.Error(w, "DB not connected", http.StatusInternalServerError)
		return
	}

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []bson.M
	if err := cursor.All(context.TODO(), &users); err != nil {
		http.Error(w, "Decode failed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h2>📄 All Subscribers</h2><ul>")
	for _, user := range users {
		fmt.Fprintf(w, "<li><b>%s</b> – %s</li>", user["name"], user["email"])
	}
	fmt.Fprint(w, "</ul><br><a href='/'>⬅️ Back</a>")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("✅ Go service healthy"))
}

func main() {
	connectMongo()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/success", successHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Go service running on port 3002")
	http.ListenAndServe(":3002", nil)
}
