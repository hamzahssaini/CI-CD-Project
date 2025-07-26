from flask import Flask, render_template_string
import requests

app = Flask(__name__)

services = {
    "Node.js": "http://node-service:3000",
    "Python": "http://python-service:3001",
    "Go": "http://go-service:3002"
}

@app.route("/")
def dashboard():
    statuses = {}
    users = []

    for name, url in services.items():
        try:
            health = requests.get(f"{url}/health", timeout=2).text
            res = requests.get(f"{url}/users", timeout=3)
            if res.status_code == 200:
                users_html = res.text
            else:
                users_html = "❌ Failed to fetch users"
            statuses[name] = {"status": health, "users": users_html}
        except Exception as e:
            statuses[name] = {"status": f"❌ {str(e)}", "users": "❌ No data"}

    return render_template_string("""
    <html>
    <head>
        <title>Microservices Dashboard</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background: #f5f5f5;
                padding: 20px;
            }
            .service {
                background: white;
                border-radius: 10px;
                padding: 20px;
                margin-bottom: 20px;
                box-shadow: 0 0 10px rgba(0,0,0,0.1);
            }
            h2 {
                color: #333;
            }
            .status {
                font-weight: bold;
                margin-bottom: 10px;
            }
        </style>
    </head>
    <body>
        <h1>📊 Microservices Dashboard</h1>
        {% for name, data in statuses.items() %}
        <div class="service">
            <h2>{{ name }}</h2>
            <p class="status">{{ data.status }}</p>
            <div>{{ data.users | safe }}</div>
        </div>
        {% endfor %}
    </body>
    </html>
    """, statuses=statuses)
