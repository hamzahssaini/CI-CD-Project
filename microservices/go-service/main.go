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

var client *mongo.Client

func healthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "✅ Go service healthy")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    mongoStatus := "❌ MongoDB not connected"
    count := int64(-1)

    if err := client.Ping(ctx, nil); err == nil {
        mongoStatus = "✅ MongoDB connected"

        collection := client.Database("myappdb").Collection("test_collection")
        cnt, err := collection.CountDocuments(ctx, bson.D{})
        if err != nil {
            log.Printf("Failed to count documents: %v", err)
        } else {
            count = cnt
        }
    }

    data := struct {
        Title          string
        MongoStatus    string
        WelcomeMessage string
        DocCount       int64
    }{
        Title:          "Go Service - Welcome",
        MongoStatus:    mongoStatus,
        WelcomeMessage: "👋 Hello from the Go Service at go.hamza.local!",
        DocCount:       count,
    }

    tmpl := template.Must(template.New("home").Parse(htmlPageWithCount))
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

const htmlPageWithCount = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>{{.Title}}</title>
<style>
  body {
    height: 100vh;
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    background: linear-gradient(135deg, #4b6cb7, #182848);
    color: #ffffff;
    text-align: center;
  }
  h1 {
    font-size: 3rem;
    margin-bottom: 0.5rem;
  }
  p {
    font-size: 1.5rem;
    margin-top: 0;
  }
  .status {
    margin-top: 1.5rem;
    font-weight: bold;
    font-size: 1.25rem;
    color: #90ee90;
  }
  .count {
    margin-top: 1rem;
    font-size: 1.1rem;
    font-style: italic;
  }
</style>
</head>
<body>
  <h1>{{.WelcomeMessage}}</h1>
  <p>Powered by Go & MongoDB</p>
  <div class="status">{{.MongoStatus}}</div>
  {{if ge .DocCount 0}}
  <div class="count">Documents in <code>test_collection</code>: {{.DocCount}}</div>
  {{else}}
  <div class="count">Unable to fetch document count.</div>
  {{end}}
</body>
</html>
`

func main() {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI environment variable not set")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    clientOpts := options.Client().ApplyURI(mongoURI)
    client, err = mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    if err = client.Ping(ctx, nil); err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }

    log.Println("Connected to MongoDB successfully!")

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/health", healthHandler)

    fmt.Println("Go service running on port 3002")
    log.Fatal(http.ListenAndServe(":3002", nil))
}
