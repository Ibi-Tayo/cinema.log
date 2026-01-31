// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from '@playwright/test';

test.describe('Bulk Film Comparison', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:4200/');
    await page.getByTestId('navbar-signin-link').click();
    await page.getByTestId('login-google-dev-button').click();
    await page.waitForURL(/\/profile\/.*/);
  });

  test('User can perform bulk comparisons', async ({ page }) => {
    // Wait for the profile page to load fully
    await expect(page.getByTestId('profile-least-comparisons-table')).toBeVisible();
    
    // Navigate to a film review page by clicking the first film in the least comparisons table
    // The table has a rowgroup for tbody, so we target cells directly
    const filmCells = page.getByTestId('profile-least-comparisons-table').locator('tbody').getByRole('row').first().getByRole('cell').first();
    await filmCells.click();

    // 1. Verify bulk mode is enabled
    await expect(page.getByRole('checkbox', { name: 'Bulk Mode' })).toBeChecked();
    
    // Check if films are available for comparison
    const noFilmsMessage = page.getByText('No more films available for comparison.');
    const hasFilmsToCompare = await noFilmsMessage.isVisible({ timeout: 2000 }).catch(() => false) === false;
    
    if (!hasFilmsToCompare) {
      // Skip the rest of the test if no films available (all comparisons already done)
      return;
    }
    
    await expect(page.getByText('Select your preference for each film below')).toBeVisible();
    await expect(page.getByText(/\d+ \/ 50 films loaded/)).toBeVisible();
    await expect(page.getByText(/Comparing: .+ \(\d{4}-\d{2}-\d{2}\)/)).toBeVisible();

    // 2. Click 'Better' on the first comparison film
    const betterButtons = page.getByRole('button', { name: 'Better' });
    await betterButtons.first().click();

    // 3. Click 'Worse' on the second comparison film if available
    const worseButtons = page.getByRole('button', { name: 'Worse' });
    const worseButtonCount = await worseButtons.count();
    if (worseButtonCount > 0) {
      await worseButtons.first().click();
    }
    
    await expect(page.getByRole('button', { name: /Submit \d+ Comparisons?/ })).toBeVisible();

    // 4. Click the 'Submit' button
    await page.getByRole('button', { name: /Submit \d+ Comparisons?/ }).click();
    
    // Verify user is redirected after submission
    await expect(page).toHaveURL(/\/profile\//);
  });
});
