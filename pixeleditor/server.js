const fs = require('fs');
const path = require('path');

const express = require('express');
const app = express();

const port = 3000;

app.get('/', (req, res) => {
    const htmlTemplate = fs.readFileSync(path.resolve(__dirname, './index.html'), 'utf8');
    res.send(htmlTemplate);
});

app.use(express.static('public'));

app.listen(port, () => console.log(`Server listening on ${port}`));
