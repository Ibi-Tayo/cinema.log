// spec: specs/basic-functions.plan.md
// seed: tests/seed.spec.ts

import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('User can sign in successfully', async ({ page }) => {
    // 1. Navigate to the home page at http://localhost:4200/
    await page.goto('http://localhost:4200/');
    await expect(page.getByRole('heading', { name: 'Your personal hub for film review' })).toBeVisible();

    // 2. Click the 'Sign In' button in the navigation bar
    await page.getByTestId('navbar-signin-link').click();
    await expect(page.getByText('Join Our Film Community')).toBeVisible();

    // 3. Click the 'Sign in with Google (Dev)' button
    await page.getByTestId('login-google-dev-button').click();
    await page.getByText('Dev Google User').first().waitFor({ state: 'visible' });
    await expect(page.url()).toContain('/profile/');
    await expect(page.getByRole('heading', { name: 'Dev Google User' })).toBeVisible();
  });
});
