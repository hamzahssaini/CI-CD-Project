const express = require('express');
const mongoose = require('mongoose');
const bodyParser = require('body-parser');
const path = require('path');

const app = express();
app.use(bodyParser.urlencoded({ extended: true }));

// MongoDB
const mongoURI = process.env.MONGO_URI;
let mongoStatus = '❌ MongoDB not connected';
let User;

mongoose.connect(mongoURI, { useNewUrlParser: true, useUnifiedTopology: true })
  .then(() => {
    mongoStatus = '✅ MongoDB connected';
    const userSchema = new mongoose.Schema({
      name: String,
      email: String,
      source: String,
    });
    User = mongoose.model('User', userSchema);
  })
  .catch(err => {
    console.error('MongoDB connection failed:', err.message);
  });

// Home Page
app.get('/', async (req, res) => {
  const userCount = User ? await User.countDocuments() : 0;

  res.send(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>Node Service - Hamza</title>
        <style>
            body {
                background: linear-gradient(135deg, #ff758c, #ff7eb3);
                font-family: Arial, sans-serif;
                color: white;
                text-align: center;
                padding: 50px;
            }
            input, button {
                padding: 10px;
                margin: 5px;
                border-radius: 5px;
                border: none;
            }
            button {
                background: white;
                color: #ff7eb3;
                font-weight: bold;
            }
        </style>
    </head>
    <body>
        <h1>👋 Welcome to Node Service at <code>node.hamza.local</code></h1>
        <p>${mongoStatus}</p>
        <p>📄 Total Subscribers: ${userCount}</p>

        <form action="/register" method="POST">
            <input name="name" placeholder="Enter Your Name" required />
            <input name="email" type="email" placeholder="Enter Your Email" required />
            <button type="submit">Subscribe</button>
        </form>

        <br>
        <a href="/users">🔍 View all users</a>
    </body>
    </html>
  `);
});

// Register Route
app.post('/register', async (req, res) => {
  if (!User) return res.status(500).send('DB error');

  const { name, email } = req.body;
  await User.create({ name, email, source: 'node-service' });
  res.redirect(`/success?name=${encodeURIComponent(name)}`);
});

// Success page
app.get('/success', (req, res) => {
  const name = req.query.name || 'User';
  res.send(`<h2>✅ Thank you, <b>${name}</b>! You're subscribed!</h2><a href="/">⬅️ Back</a>`);
});

// Users Route
app.get('/users', async (req, res) => {
  if (!User) return res.status(500).send('DB not connected');

  const users = await User.find({});
  const list = users.map(u => `<li>${u.name} (${u.email})</li>`).join('');
  res.send(`<ul>${list}</ul><a href="/">⬅️ Back</a>`);
});

// Health check
app.get('/health', (req, res) => {
  res.send('✅ Node service healthy');
});

// Start server
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Node.js service listening on port ${PORT}`);
});
