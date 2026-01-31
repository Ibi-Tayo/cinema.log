// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:4200/');
    await page.getByTestId('navbar-signin-link').click();
    await page.getByTestId('login-google-dev-button').click();
    await page.waitForURL(/\/profile\/.*/);
  });

  test('Unauthenticated user cannot access protected pages', async ({ page }) => {
    // First sign out to test protected routes
    await page.getByTestId('navbar-profile-link').click();
    await page.getByRole('menuitem', { name: 'Sign Out' }).click();

    // 1. Navigate directly to a profile URL without being logged in
    await page.goto('http://localhost:4200/profile/c55bda9d-5434-4df6-81c2-52a9e0d63c55');
    await expect(page.getByText('Join Our Film Community')).toBeVisible();

    // 2. Navigate directly to the search page without being logged in
    await page.goto('http://localhost:4200/search');
    await expect(page.getByText('Join Our Film Community')).toBeVisible();

    // 3. Navigate directly to the recommendations page without being logged in
    await page.goto('http://localhost:4200/recommendations/c55bda9d-5434-4df6-81c2-52a9e0d63c55');
    await expect(page.getByText('Join Our Film Community')).toBeVisible();
  });
});
