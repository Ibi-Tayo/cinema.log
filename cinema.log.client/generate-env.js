/**
 * Generates runtime environment configuration for the Angular app.
 *
 * This script creates src/assets/env.js which is loaded by index.html before the app starts.
 * It allows environment variables to be set at build time (e.g., on Railway) rather than
 * being hardcoded in the TypeScript environment files.
 *
 * Environment Variables:
 * - API_URL: Backend API URL (e.g., https://reverse-proxy-cinemalog-pr-123.up.railway.app/api)
 * - AUTH_DOMAIN: Authentication domain (optional)
 * - ENVIRONMENT: Environment name (development, staging, production)
 *
 * Usage:
 * - This script runs automatically before `npm run build` and `npm test`
 * - On Railway, set the API_URL environment variable in your service configuration
 * - For PR deployments, use: https://reverse-proxy-cinemalog-pr-${{ github.event.pull_request.number }}.up.railway.app/api
 */

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
