import { test as setup, expect } from "@playwright/test";

const authFile = ".auth/user.json";
const AUTH_TIMEOUT = 10000;

setup("authenticate with dev login", async ({ page, context, baseURL }) => {
  // In CI/test environments, use the dev login endpoint to bypass OAuth
  // This creates a test user automatically without requiring GitHub credentials

  // Navigate to the home page to verify authentication
  await page.goto("/");
  await expect(
    page.getByRole("heading", { name: "Your personal hub for film review" }),
  ).toBeVisible();

  await page.getByTestId("navbar-signin-link").click();

  await page.getByTestId("login-github-dev-button").click();

  await page.reload();

  // Verify user is logged in by checking for profile link in navigation
  // The dev user should be automatically created and logged in
  await expect(page.getByTestId("navbar-profile-link")).toBeVisible({
    timeout: AUTH_TIMEOUT,
  });

  // Save authentication state (cookies) for reuse in tests
  await context.storageState({ path: authFile });
});
