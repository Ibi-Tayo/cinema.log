// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Review and Rating", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("User can submit a complete film review", async ({ page }) => {
    // Navigate to search page to find a film to review
    await page.getByTestId("navbar-search-link").click();

    // 1. Search for and select a film (e.g., 'Casablanca')
    await page.getByTestId("search-input").fill("Casablanca");
    await page.getByTestId("search-input").press("Enter");

    await page.getByRole("img", { name: "Casablanca" }).first().click();
    await expect(
      page.getByRole("heading", { name: "Casablanca" }),
    ).toBeVisible();

    // 2. Check if this is a new review (star buttons visible) or an update
    const starButton = page.getByTestId("review-form-star-5");
    const isNewReview = await starButton
      .isVisible({ timeout: 2000 })
      .catch(() => false);

    if (isNewReview) {
      // New review flow - click the 5-star rating button
      await starButton.click();
      await expect(page.getByText("Click to rate (5)")).toBeVisible();

      // 3. Type a review in the 'Your Thoughts' textbox
      await page
        .getByTestId("review-form-textarea")
        .fill("A masterpiece of storytelling with unforgettable dialogue");

      // 4. Click the 'Submit Review' button
      await page.getByTestId("review-form-submit-button").click();
      await expect(
        page.getByText("Review submitted successfully"),
      ).toBeVisible();
      await expect(
        page.getByRole("heading", { name: "Update Review" }),
      ).toBeVisible();
    } else {
      // Update review flow - film has already been reviewed
      await expect(
        page.getByRole("heading", { name: "Update Review" }),
      ).toBeVisible();

      // Update the review text
      await page
        .getByTestId("review-form-textarea")
        .fill("Updated review: A timeless classic!");
      await page.getByRole("button", { name: "Update Review" }).click();
      await expect(
        page.getByText(/Review (updated|submitted) successfully/i),
      ).toBeVisible();
    }
  });
});
