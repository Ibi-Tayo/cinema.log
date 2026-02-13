// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Film Search", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("User can search for films by title", async ({ page }) => {
    // 1. Click 'Search' in the navigation bar
    await page.getByTestId("navbar-search-link").click();
    await expect(
      page.getByRole("heading", { name: "Search for Films" }),
    ).toBeVisible();
    await expect(
      page.getByRole("textbox", { name: "Search for a film by title..." }),
    ).toBeVisible();
    await expect(page.getByText("Start Your Search")).toBeVisible();

    // 2. Type 'The Dark Knight' in the search textbox
    await page.getByTestId("search-input").fill("The Dark Knight");

    // 3. Press Enter or wait for search results
    await page.getByTestId("search-input").press("Enter");
    await expect(
      page.getByRole("img", { name: "The Dark Knight" }).first(),
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { name: "The Dark Knight Rises" }).first(),
    ).toBeVisible();
  });
});
