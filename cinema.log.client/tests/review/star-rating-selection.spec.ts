// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Review and Rating", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("User can select different star ratings", async ({ page }) => {
    // Navigate to search page
    await page.getByTestId("navbar-search-link").click();
    await page.waitForURL("**/search");

    // 1. Navigate to a film review page
    await page.getByTestId("search-input").fill("Schindler's List");
    await page.getByTestId("search-input").press("Enter");

    await page
      .getByRole("img", { name: "Schindler's List", exact: true })
      .click();

    const secondStar = page.locator('span.star').nth(2); // 1-based index for stars because there is a hidden 0-star element
    await secondStar.click();

    const ratingHint = page.locator('p.rating-hint');
    // expect to have either 1.5 or 2 stars selected
    await expect(ratingHint).toHaveText(/\s*Click to rate \(1\.5|2\)\s*/);
  });
});
