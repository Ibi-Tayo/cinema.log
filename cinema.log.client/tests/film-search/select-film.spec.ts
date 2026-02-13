// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Search", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("User can click on a search result to review", async ({ page }) => {
    // Navigate to search page
    await page.getByTestId("navbar-search-link").click();
    await page.waitForURL("**/search");

    // 1. Search for 'Inception'
    await page.getByTestId("search-input").fill("Inception");
    await page.getByTestId("search-input").press("Enter");
 
    await expect(
      page.getByRole("img", { name: "Inception" }).first(),
    ).toBeVisible();

    // 2. Click on the 'Inception' film card from search results
    await page.getByRole("img", { name: "Inception" }).first().click();
    await expect(page.url()).toContain("/review/");
    await expect(
      page.getByRole("heading", { name: "Inception" }),
    ).toBeVisible();
    await expect(page.getByTestId("film-display-year")).toBeVisible();
  });
});
