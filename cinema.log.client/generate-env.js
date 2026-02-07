const fs = require("fs");
const path = require("path");

const envConfig = `
(function(window) {
  window.__env = window.__env || {};
  window.__env.apiUrl = '${process.env.API_URL || "http://localhost:8080"}';
  window.__env.authDomain = '${process.env.AUTH_DOMAIN || ""}';
  window.__env.environment = '${process.env.ENVIRONMENT || "development"}';
}(this));
`;

const targetPath = path.join(__dirname, "src/assets/env.js");
fs.mkdirSync(path.dirname(targetPath), { recursive: true });
fs.writeFileSync(targetPath, envConfig);

console.log("Environment configuration generated at src/assets/env.js");
