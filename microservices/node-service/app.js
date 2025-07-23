const express = require('express');
const app = express();

app.get('/health', (req, res) => {
  res.status(200).send('✅ Node service healthy');
});

app.get('/', (req, res) => {
  res.send('👋 Hello from Node.js Service!');
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Node.js service listening on port ${PORT}`);
});
