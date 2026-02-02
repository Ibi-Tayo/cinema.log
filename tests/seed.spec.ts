import { test, expect } from "@playwright/test";
import { ensureFilmExists } from "./utils/test-helpers";

test.describe("Test group", () => {
  test("seed", async ({ page }) => {
    // This test verifies authentication is working
    // Authentication state is already handled by auth.setup.ts
    await page.goto("/");

    // Verify user is authenticated
    await expect(page.getByTestId("navbar-profile-link")).toBeVisible();

    const title = page.locator("h1");
    await expect(title).toHaveText("Your personal hub for film review");
  });
});
