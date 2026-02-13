// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Authentication", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("Unauthenticated user cannot access protected pages", async ({
    page,
  }) => {
    // First sign out to test protected routes
    await page.getByTestId("navbar-profile-link").click();
    await page.getByRole("menuitem", { name: "Sign Out" }).click();

    // Wait for redirect to login/home page after sign out
    await page.waitForURL(/\/(login|home)/);

    // 1. Navigate directly to a profile URL without being logged in
    await page.goto("/profile/test-user-id");
    await expect(page.getByText("Join Our Film Community")).toBeVisible();

    // 2. Navigate directly to the search page without being logged in
    await page.goto("/search");
    await expect(page.getByText("Join Our Film Community")).toBeVisible();

    // 3. Navigate directly to the recommendations page without being logged in
    await page.goto("/recommendations/test-user-id");
    await expect(page.getByText("Join Our Film Community")).toBeVisible();
  });
});
