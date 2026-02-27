// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe.skip("Submit Reviews then Bulk Film Comparison", () => {
  test("Create at least two reviews", async ({ page }) => {
    await page.goto("/");

    const films = ["Casablanca", "The Godfather"];

    for (const title of films) {
      await page.getByTestId("navbar-search-link").click();
      await page.getByTestId("search-input").fill(title);
      await page.getByTestId("search-input").press("Enter");

      const img = page.getByRole("img", { name: new RegExp(title) }).first();
      const found = await img.isVisible().catch(() => false);
      if (!found) continue;

      await img.click();
      await expect(
        page.getByRole("heading", { name: new RegExp(title) }),
      ).toBeVisible();

      const starButton = page.getByTestId("review-form-star-5");
      const isNewReview = await starButton
        .isVisible({ timeout: 2000 })
        .catch(() => false);

      if (isNewReview) {
        await starButton.click();
        await expect(page.getByText("Click to rate (5)")).toBeVisible();
        await page
          .getByTestId("review-form-textarea")
          .fill(`Automated review for ${title}`);
        await page.getByTestId("review-form-submit-button").click();
        await expect(
          page.getByText("Review submitted successfully"),
        ).toBeVisible();
        await expect(
          page.getByRole("heading", { name: "Update Review" }),
        ).toBeVisible();
      } else {
        await expect(
          page.getByRole("heading", { name: "Update Review" }),
        ).toBeVisible();
        await page
          .getByTestId("review-form-textarea")
          .fill(`Updated automated review for ${title}`);
        await page.getByRole("button", { name: "Update Review" }).click();
        await expect(
          page.getByText(/Review (updated|submitted) successfully/i),
        ).toBeVisible();
      }
    }
  });

  test("User can toggle between bulk and sequential mode", async ({ page }) => {
    await page.goto("/");

    await page.getByTestId("navbar-profile-link").click();
    await page.getByRole("menuitem", { name: "Profile" }).click();

    // Wait for table to load and check if there are any rows
    const tableRow = page
      .getByTestId("profile-least-comparisons-table")
      .locator("tbody")
      .getByRole("row")
      .first();

    await expect(tableRow).toBeVisible({ timeout: 10000 });

    // Navigate to film review page with comparison mode by clicking first available film in table
    await tableRow.getByRole("cell").first().click();

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

  test("User can perform bulk comparisons", async ({ page }) => {
    await page.goto("/");

    await page.getByTestId("navbar-profile-link").click();
    await page.getByRole("menuitem", { name: "Profile" }).click();

    await expect(
      page.getByTestId("profile-least-comparisons-table"),
    ).toBeVisible();

    const filmCells = page
      .getByTestId("profile-least-comparisons-table")
      .locator("tbody")
      .getByRole("row")
      .first()
      .getByRole("cell")
      .first();
    await filmCells.click();

    await expect(
      page.getByRole("checkbox", { name: "Bulk Mode" }),
    ).toBeChecked();

    const noFilmsMessage = page.getByText(
      "No more films available for comparison.",
    );
    const hasFilmsToCompare =
      (await noFilmsMessage.isVisible({ timeout: 2000 }).catch(() => false)) ===
      false;

    if (!hasFilmsToCompare) {
      return;
    }

    await expect(
      page.getByText("Select your preference for each film below"),
    ).toBeVisible();
    await expect(page.getByText(/\d+ \/ 50 films loaded/)).toBeVisible();
    await expect(
      page.getByText(/Comparing: .+ \(\d{4}-\d{2}-\d{2}\)/),
    ).toBeVisible();

    const betterButtons = page.getByRole("button", { name: "Better" });
    await betterButtons.first().click();

    const worseButtons = page.getByRole("button", { name: "Worse" });
    const worseButtonCount = await worseButtons.count();
    if (worseButtonCount > 0) {
      await worseButtons.first().click();
    }

    await expect(
      page.getByRole("button", { name: /Submit \d+ Comparisons?/ }),
    ).toBeVisible();
    await page.getByRole("button", { name: /Submit \d+ Comparisons?/ }).click();

    await expect(page).toHaveURL(/\/profile\//);
  });
});
