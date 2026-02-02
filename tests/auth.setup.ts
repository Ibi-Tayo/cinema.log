import { test as setup, expect } from "@playwright/test";
import * as fs from "fs";
import * as path from "path";

const authFile = ".auth/user.json";

setup("authenticate with GitHub OAuth", async ({ page, context }) => {
  const testEmail = process.env.TEST_GITHUB_EMAIL;
  const testPassword = process.env.TEST_GITHUB_PASSWORD;

  if (!testEmail || !testPassword) {
    throw new Error(
      "TEST_GITHUB_EMAIL and TEST_GITHUB_PASSWORD environment variables must be set",
    );
  }

  // Navigate to the home page
  await page.goto("/");
  await expect(
    page.getByRole("heading", { name: "Your personal hub for film review" }),
  ).toBeVisible();

  // Click Sign In
  await page.getByTestId("navbar-signin-link").click();
  await expect(page.getByText("Join Our Film Community")).toBeVisible();

  // Click Sign in with GitHub
  await page.getByTestId("login-github-button").click();

  // Wait for GitHub login page
  await page.waitForURL(/github\.com\/login/);

  // Fill in GitHub credentials
  await page.fill('input[name="login"]', testEmail);
  await page.fill('input[name="password"]', testPassword);
  await page.click('input[type="submit"]');

  // Handle authorization if needed (first time)
  const authorizeButton = page.locator('button[name="authorize"]');
  if (await authorizeButton.isVisible({ timeout: 5000 }).catch(() => false)) {
    await authorizeButton.click();
  }

  // Wait for redirect back to the app with profile URL
  await page.waitForURL(/\/profile\/.*/);

  // Verify user is logged in by checking for user name in navigation
  await expect(page.getByTestId("navbar-profile-link")).toBeVisible();

  // Ensure .auth directory exists
  const authDir = path.dirname(authFile);
  if (!fs.existsSync(authDir)) {
    fs.mkdirSync(authDir, { recursive: true });
  }

  // Save authentication state
  await context.storageState({ path: authFile });
});
