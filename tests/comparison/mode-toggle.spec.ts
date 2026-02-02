// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Bulk Film Comparison", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");
  });

  test("User can toggle between bulk and sequential mode", async ({ page }) => {
    // Navigate to film review page with comparison mode
    await page
      .getByTestId("profile-least-comparisons-table")
      .getByRole("cell", { name: "The Matrix" })
      .click();

    // 1. Verify bulk mode is active
    await expect(
      page.getByRole("checkbox", { name: "Bulk Mode" }),
    ).toBeChecked();
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
