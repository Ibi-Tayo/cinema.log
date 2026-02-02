// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Review and Rating", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");
  });

  test("User can select different star ratings", async ({ page }) => {
    // Navigate to search page
    await page.getByTestId("navbar-search-link").click();

    // 1. Navigate to a film review page
    await page.getByTestId("search-input").fill("Schindler's List");
    await page.getByTestId("search-input").press("Enter");
    await page
      .getByRole("img", { name: "Schindler's List", exact: true })
      .click();
    await expect(
      page.getByRole("button", { name: "Rate 1 out of 5 stars" }),
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Rate 5 out of 5 stars" }),
    ).toBeVisible();

    // 2. Click the 3-star button
    await page.getByTestId("review-form-star-3").click();
    await expect(page.getByText("Click to rate (3)")).toBeVisible();

    // 3. Click the 5-star button
    await page.getByTestId("review-form-star-5").click();
    await expect(page.getByText("Click to rate (5)")).toBeVisible();

    // 4. Click the 1-star button
    await page.getByTestId("review-form-star-1").click();
    await expect(page.getByText("Click to rate (1)")).toBeVisible();
  });
});
