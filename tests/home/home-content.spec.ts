// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from "@playwright/test";

test.describe("Home Page", () => {
  test.beforeEach(async ({ page }) => {
    // Authentication is already handled by auth.setup.ts
    await page.goto("/");
  });

  test("Home page displays correct content", async ({ page }) => {
    // 1. Navigate to the home page
    await page.getByTestId("navbar-logo").click();

    // expect: The main heading 'Your personal hub for film review' is displayed
    await expect(
      page.getByRole("heading", { name: "Your personal hub for film review" }),
    ).toBeVisible();

    // expect: A subheading 'A smarter way to track and rate the films you love' is shown
    await expect(
      page.getByText("A smarter way to track and rate the films you love"),
    ).toBeVisible();

    // expect: Three main content sections are visible: 'Review & Rate Films', 'ELO Rating System', and 'Ready to Start?'
    await expect(
      page.getByRole("heading", { name: "Review & Rate Films" }),
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { name: "ELO Rating System" }),
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { name: "Ready to Start?" }),
    ).toBeVisible();

    // 2. Scroll to view the 'ELO Rating System' section
    // expect: A 'How It Works' explanation with three numbered steps is displayed
    await expect(
      page.getByRole("heading", { name: "How It Works" }),
    ).toBeVisible();

    // expect: Step 1: 'Compare Films' with description
    await expect(page.getByText("Compare Films:")).toBeVisible();

    // expect: Step 2: 'Dynamic Ratings' with description
    await expect(page.getByText("Dynamic Ratings:")).toBeVisible();

    // expect: Step 3: 'Refined Over Time' with description
    await expect(page.getByText("Refined Over Time:")).toBeVisible();

    // 3. Review the 'Review & Rate Films' section
    // expect: Text explains the 1-5 star rating system
    await expect(
      page.getByText("Rate films from 1 to 5 stars based on your experience"),
    ).toBeVisible();
  });
});
