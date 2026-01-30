import { test, expect } from "@playwright/test";

test.describe("Test group", () => {
  test("seed", async ({ page }) => {
    await page.goto("http://localhost:4200/");
    const title = page.locator("h1");
    await expect(title).toHaveText("Your personal hub for film review");
    // click sign in button in nav
    await page.click("text=Sign In");
    // click button to sign in with google  Sign in with Google (Dev)
    await page.click("text=Sign in with Google (Dev)");
    // wait for redirect to http://localhost:4200/profile/c55bda9d-5434-4df6-81c2-52a9e0d63c55
    await page.waitForURL(/\/profile\/.*/);
    // check user is signed in because sign in nav is replaced with Dev Google User
    const user = page.getByRole("navigation").getByText("Dev Google User");
    await expect(user).toBeVisible();
  });
});
