// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Review and Rating", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");
  });

  test("User cannot submit review without rating", async ({ page }) => {
    // Navigate to search page
    await page.getByTestId("navbar-search-link").click();

    // 1. Navigate to review page for a new film
    await page.getByTestId("search-input").fill("Forrest Gump");
    await page.getByTestId("search-input").press("Enter");
    await expect(
      page.getByRole("heading", { name: /Search Results/ }),
    ).toBeVisible();
    await page.getByRole("img", { name: "Forrest Gump" }).first().click();
    await expect(
      page.getByRole("heading", { name: "Review Film" }),
    ).toBeVisible();

    // 2. Type review text without selecting a star rating
    await page
      .getByTestId("review-form-textarea")
      .fill("This is a great movie but I haven't rated it");
    await expect(page.getByText("Click to rate (Not rated)")).toBeVisible();

    // 3. Attempt to click the disabled 'Submit Review' button
    await expect(page.getByTestId("review-form-submit-button")).toBeDisabled();
  });
});
