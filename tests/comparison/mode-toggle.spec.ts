// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Bulk Film Comparison", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");

    // Navigate to profile page
    await page.getByTestId("navbar-profile-link").click();
    await page.getByRole("menuitem", { name: "Profile" }).click();
  });

  test("User can toggle between bulk and sequential mode", async ({ page }) => {
    // Navigate to film review page with comparison mode by clicking first available film in table
    await page
      .getByTestId("profile-least-comparisons-table")
      .locator("tbody")
      .getByRole("row")
      .first()
      .getByRole("cell")
      .first()
      .click();

    // 1. Verify bulk mode is active
    await expect(
      page.getByRole("checkbox", { name: "Bulk Mode" }),
    ).toBeChecked();

    // Check if films are available for comparison
    const noFilmsMessage = page.getByText(
      "No more films available for comparison.",
    );
    const noFilmsVisible = await noFilmsMessage
      .isVisible({ timeout: 2000 })
      .catch(() => false);

    if (noFilmsVisible) {
      // Skip the rest of the test if no films available (all comparisons already done)
      return;
    }

    await expect(
      page.getByText("Select your preference for each film below"),
    ).toBeVisible();

    // 2. Click the 'Bulk Mode' checkbox to uncheck it
    await page.getByTestId("review-mode-toggle-input").click();
    await expect(page.getByText(/Comparison \d+ \/ \d+/)).toBeVisible();
    await expect(page.getByText("Which film do you prefer?")).toBeVisible();

    // 3. Click the checkbox again to re-enable bulk mode
    await page.getByTestId("review-mode-toggle-input").click();
    await expect(
      page.getByText("Select your preference for each film below"),
    ).toBeVisible();
  });
});
