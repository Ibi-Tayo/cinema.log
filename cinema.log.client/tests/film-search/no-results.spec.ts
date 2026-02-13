// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Search", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("Search displays empty state for no results", async ({ page }) => {
    // 1. Navigate to the search page
    await page.getByTestId("navbar-search-link").click();
    await expect(
      page.getByRole("heading", { name: "Search for Films" }),
    ).toBeVisible();

    // 2. Type a nonsensical string that will return no results
    await page.getByTestId("search-input").fill("xyzabc123nonexistentfilm");
    await page.getByTestId("search-input").press("Enter");

    // 3. Wait for search completion
    await expect(
      page.getByText("No films found matching your search."),
    ).toBeVisible();
  });
});
