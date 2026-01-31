// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from '@playwright/test';

test.describe('Film Review and Rating', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:4200/');
    await page.getByTestId('navbar-signin-link').click();
    await page.getByTestId('login-google-dev-button').click();
    await page.waitForURL(/\/profile\/.*/);
  });

  test('User can navigate back to search from review page', async ({ page }) => {
    // Navigate to search page
    await page.getByTestId('navbar-search-link').click();

    // 1. Search for and select a film to review
    await page.getByTestId('search-input').fill('Gladiator');
    await page.getByTestId('search-input').press('Enter');
    await expect(page.getByRole('heading', { name: /Search Results/ })).toBeVisible();
    await page.getByRole('img', { name: 'Gladiator' }).first().click();
    await expect(page.getByRole('button', { name: 'Back to Search' })).toBeVisible();

    // 2. Click the 'Back to Search' button
    await page.getByTestId('review-back-button').click();
    await expect(page.getByRole('heading', { name: 'Search for Films' })).toBeVisible();
  });
});
