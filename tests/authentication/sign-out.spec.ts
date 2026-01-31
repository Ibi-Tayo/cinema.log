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

  test('User can sign out successfully', async ({ page }) => {
    // 1. Click the 'Dev Google User' dropdown in the navigation bar
    await page.getByTestId('navbar-profile-link').click();
    await expect(page.getByRole('menuitem', { name: 'Profile' })).toBeVisible();
    await expect(page.getByRole('menuitem', { name: 'Sign Out' })).toBeVisible();

    // 2. Click 'Sign Out' from the dropdown menu
    await page.getByRole('menuitem', { name: 'Sign Out' }).click();
    await expect(page.getByRole('link', { name: 'Sign in' })).toBeVisible();
  });
});
