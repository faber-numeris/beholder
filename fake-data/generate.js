// app.js
const express = require("express");
const { faker } = require("@faker-js/faker");

const app = express();

app.get("/user", (req, res) => {
  res.json({
    firstName: faker.person.firstName(),
    lastName: faker.person.lastName(),
    email: faker.internet.email(),
  });
});

const PORT = 3001;
app.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});