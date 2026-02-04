import { test as setup, expect } from "@playwright/test";
import * as fs from "fs";
import * as path from "path";

const authFile = ".auth/user.json";

setup("authenticate with dev login", async ({ page, context, baseURL }) => {
  // In CI/test environments, use the dev login endpoint to bypass OAuth
  // This creates a test user automatically without requiring GitHub credentials
  
  // Call the dev login endpoint directly to get authentication cookies
  const response = await page.request.get(`${baseURL}/api/auth/dev/login`, {
    failOnStatusCode: false,
  });

  if (!response.ok()) {
    throw new Error(
      `Dev login failed with status ${response.status()}. ` +
      `This endpoint is only available in non-production environments.`
    );
  }

  // Extract cookies from the response
  const cookies = await context.cookies();
  
  // Navigate to the home page to verify authentication
  await page.goto("/");
  await expect(
    page.getByRole("heading", { name: "Your personal hub for film review" }),
  ).toBeVisible();

  // Verify user is logged in by checking for profile link in navigation
  // The dev user should be automatically created and logged in
  await expect(page.getByTestId("navbar-profile-link")).toBeVisible({ timeout: 10000 });

  // Ensure .auth directory exists
  const authDir = path.dirname(authFile);
  if (!fs.existsSync(authDir)) {
    fs.mkdirSync(authDir, { recursive: true });
  }

  // Save authentication state (cookies) for reuse in tests
  await context.storageState({ path: authFile });
});

