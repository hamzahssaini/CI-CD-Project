from flask import Flask, request, redirect
from pymongo import MongoClient
import os

app = Flask(__name__)

# MongoDB setup
mongo_uri = os.getenv("MONGO_URI")
try:
    client = MongoClient(mongo_uri, serverSelectionTimeoutMS=3000)
    db = client["myappdb"]
    collection = db["users"]
    mongo_status = "✅ MongoDB connected"
except Exception as e:
    collection = None
    mongo_status = f"❌ MongoDB connection failed: {e}"

@app.route("/")
def home():
    user_count = collection.count_documents({}) if collection else 0
    return f"""
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8" />
        <title>Python Service – Hamza</title>
        <style>
            body {{
                font-family: Arial, sans-serif;
                background: linear-gradient(120deg, #00c6ff, #0072ff);
                color: white;
                text-align: center;
                padding: 50px;
            }}
            h1 {{
                font-size: 2.5rem;
            }}
            form {{
                margin-top: 30px;
            }}
            input {{
                padding: 10px;
                margin: 5px;
                border-radius: 5px;
                border: none;
            }}
            button {{
                padding: 10px 20px;
                border-radius: 5px;
                background-color: #fff;
                color: #0072ff;
                font-weight: bold;
                border: none;
                cursor: pointer;
            }}
            a {{
                color: #fff;
                text-decoration: underline;
            }}
        </style>
    </head>
    <body>
        <h1>👋 Welcome to Python Service at <code>python.hamza.local</code></h1>
        <p>{mongo_status}</p>
        <p>📄 Total Subscribers: <strong>{user_count}</strong></p>

        <form action="/register" method="POST">
            <input name="name" placeholder="Enter Your Name" required />
            <input name="email" type="email" placeholder="Enter Your Email" required />
            <button type="submit">Subscribe</button>
        </form>

        <br/>
        <a href="/users">🔍 View all subscribers</a>
    </body>
    </html>
    """

@app.route("/register", methods=["POST"])
def register():
    name = request.form.get("name")
    email = request.form.get("email")
    if collection:
        collection.insert_one({
            "name": name,
            "email": email,
            "source": "python-service"
        })
        return redirect(f"/success?name={name}")
    return "❌ DB error", 500

@app.route("/success")
def success():
    name = request.args.get("name", "User")
    return f"""
    <h2>✅ Thank you, <b>{name}</b>! You're subscribed successfully.</h2>
    <a href="/">⬅️ Back to Home</a>
    """

@app.route("/users")
def list_users():
    if not collection:
        return "❌ DB error", 500
    users = list(collection.find({}, {"_id": 0}))
    return {
        "count": len(users),
        "users": users
    }

@app.route("/health")
def health():
    return "✅ Python service healthy", 200

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=3001)
