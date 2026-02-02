// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("User Profile", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");
  });

  test("Profile displays user information correctly", async ({ page }) => {
    // 1. Click the user dropdown and select 'Profile'
    await page.getByTestId("navbar-profile-link").click();
    await page.getByRole("menuitem", { name: "Profile" }).click();

    // Verify profile page elements
    await expect(
      page.getByRole("heading", { name: "Dev Google User" }),
    ).toBeVisible();
    await expect(page.getByText("@devgoogleuser")).toBeVisible();
    await expect(page.getByText("Member since January 2026")).toBeVisible();
    await expect(
      page.getByRole("img", { name: "Dev Google User" }),
    ).toBeVisible();
  });
});
