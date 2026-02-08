// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Review and Rating", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("User can navigate back to search from review page", async ({
    page,
  }) => {
    // Navigate to search page
    await page.getByTestId("navbar-search-link").click();
    await page.waitForURL("**/search");

    // 1. Search for and select a film to review
    await page.getByTestId("search-input").fill("Gladiator");
    await page.getByTestId("search-input").press("Enter");

    await page.getByRole("img", { name: "Gladiator" }).first().click();
    await expect(
      page.getByRole("button", { name: "Back to Search" }),
    ).toBeVisible();

    // 2. Click the 'Back to Search' button
    await page.getByTestId("review-back-button").click();
    await expect(
      page.getByRole("heading", { name: "Search for Films" }),
    ).toBeVisible();
  });
});
