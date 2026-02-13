import { Page, expect } from "@playwright/test";

/**
 * Ensures a film exists in the user's collection by searching for it and adding it if not present.
 * This makes tests idempotent by creating required test data on demand.
 *
 * @param page - The Playwright page object
 * @param filmTitle - The title of the film to ensure exists
 * @param rating - Optional star rating (1-5) to give the film
 * @param reviewText - Optional review text
 */
export async function ensureFilmExists(
  page: Page,
  filmTitle: string,
  rating?: number,
  reviewText?: string,
): Promise<void> {
  // Navigate to search page
  await page.goto("/search");
  await expect(
    page.getByRole("heading", { name: "Search for Films" }),
  ).toBeVisible();

  // Search for the film
  await page.getByTestId("search-input").fill(filmTitle);
  await page.getByTestId("search-input").press("Enter");

  // Wait for search results
  await expect(
    page.getByRole("heading", { name: /Search Results/ }),
  ).toBeVisible();

  // Click on the first matching film
  const filmCard = page
    .locator('.film-card, [data-testid^="film-card"]')
    .first();
  await filmCard.click();

  // Check if already reviewed
  const alreadyReviewed = await page
    .getByText("You have already reviewed this film")
    .isVisible({ timeout: 2000 })
    .catch(() => false);

  if (alreadyReviewed) {
    // Film already exists in user's collection
    return;
  }

  // Add review if not already reviewed
  if (rating) {
    // Click the star rating
    const stars = page
      .locator('.star-rating .star, [data-testid="star-rating"] button')
      .nth(rating - 1);
    await stars.click();
  }

  if (reviewText) {
    const reviewTextarea = page.locator(
      'textarea[placeholder*="review"], textarea[name="review"]',
    );
    await reviewTextarea.fill(reviewText);
  }

  // Submit the review
  const submitButton = page.getByRole("button", { name: /submit|add review/i });
  await submitButton.click();

  // Wait for success confirmation
  await expect(
    page
      .getByText(/review submitted|successfully added/i)
      .or(page.locator(".success-message")),
  ).toBeVisible({ timeout: 5000 });
}

/**
 * Ensures the user is authenticated. If not, performs re-authentication.
 * This provides a fallback when the session expires during test execution.
 *
 * @param page - The Playwright page object
 */
export async function ensureAuthenticated(page: Page): Promise<void> {
  // Navigate to home page to check authentication status
  await page.goto("/");

  // Wait for page to load
  await page.waitForLoadState("networkidle");

  // Check if user is authenticated by looking for user element in navbar
  const userInNav = await page
    .getByTestId("navbar-profile-link")
    .isVisible({ timeout: 2000 })
    .catch(() => false);

  if (!userInNav) {
    // Not authenticated, need to sign in
    const testEmail = process.env.TEST_GITHUB_EMAIL;
    const testPassword = process.env.TEST_GITHUB_PASSWORD;

    if (!testEmail || !testPassword) {
      throw new Error(
        "TEST_GITHUB_EMAIL and TEST_GITHUB_PASSWORD environment variables must be set for re-authentication",
      );
    }

    // Click Sign in button
    await page.getByTestId("navbar-signin-link").click();
    await expect(page.getByText("Join Our Film Community")).toBeVisible();

    // Click Sign in with GitHub
    await page.getByTestId("login-github-button").click();

    // Wait for GitHub login page
    await page.waitForURL(/github\.com\/login/);

    // Fill in GitHub credentials
    await page.fill('input[name="login"]', testEmail);
    await page.fill('input[name="password"]', testPassword);
    await page.click('input[type="submit"]');

    // Handle authorization if needed
    const authorizeButton = page.locator('button[name="authorize"]');
    if (await authorizeButton.isVisible({ timeout: 5000 }).catch(() => false)) {
      await authorizeButton.click();
    }

    // Wait for redirect back to the app
    await page.waitForURL(/\/profile\/.*/);
  }

  // Verify authentication by checking for user in navigation
  await expect(page.getByTestId("navbar-profile-link")).toBeVisible();
}

/**
 * Signs in using the dev login (only works in non-production environments)
 * This is kept for local testing compatibility.
 *
 * @param page - The Playwright page object
 * @deprecated Use ensureAuthenticated() instead for Railway PR environments
 */
export async function devLogin(page: Page): Promise<void> {
  await page.goto("/");
  await page.getByTestId("navbar-signin-link").click();
  await page.getByTestId("login-google-dev-button").click();
  await page.waitForURL(/\/profile\/.*/);
}
