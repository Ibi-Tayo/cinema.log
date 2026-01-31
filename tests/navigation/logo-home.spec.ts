// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:4200/');
    await page.getByTestId('navbar-signin-link').click();
    await page.getByTestId('login-google-dev-button').click();
    await page.waitForURL(/\/profile\/.*/);
  });

  test('Logo navigation returns user to home', async ({ page }) => {
    // 1. Navigate to any page in the application (search page)
    await page.getByTestId('navbar-search-link').click();
    await expect(page.url()).toContain('/search');

    // 2. Click the '_cinema.log()' logo in the navigation bar
    await page.getByTestId('navbar-logo').click();
    await expect(page.url()).toContain('/home');
    await expect(page.getByRole('heading', { name: 'Your personal hub for film review' })).toBeVisible();
  });
});
