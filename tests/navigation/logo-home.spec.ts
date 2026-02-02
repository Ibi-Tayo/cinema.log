// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Navigation", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");
  });

  test("Logo navigation returns user to home", async ({ page }) => {
    // 1. Navigate to any page in the application (search page)
    await page.getByTestId("navbar-search-link").click();
    await expect(page.url()).toContain("/search");

    // 2. Click the '_cinema.log()' logo in the navigation bar
    await page.getByTestId("navbar-logo").click();
    await expect(page.url()).toContain("/home");
    await expect(
      page.getByRole("heading", { name: "Your personal hub for film review" }),
    ).toBeVisible();
  });
});
