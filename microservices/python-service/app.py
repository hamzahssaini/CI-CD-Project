from flask import Flask
app = Flask(__name__)

@app.route('/health')
def health():
    return "✅ Python service healthy", 200

@app.route('/')
def home():
    return "👋 Hello from Python Service!"

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=3001)
