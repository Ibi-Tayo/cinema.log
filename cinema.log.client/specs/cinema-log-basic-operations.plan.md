# Cinema.log Basic Operations Test Plan

## Application Overview

A comprehensive test plan for Cinema.log, a film rating application that uses an ELO rating system to help users build personalized film rankings. The application allows users to search for films, write reviews with star ratings, compare films to generate ELO rankings, and view their profile with film statistics and rankings.

## Test Scenarios

### 1. Authentication

**Seed:** `seed.spec.ts`

#### 1.1. User can sign in successfully

**File:** `tests/authentication/sign-in.spec.ts`

**Steps:**
  1. Navigate to the home page at 'http://localhost:4200/'
    - expect: The home page loads with the heading 'Your personal hub for film review'
  2. Click the 'Sign In' button in the navigation bar
    - expect: The sign-in page or modal appears
  3. Click the 'Sign in with Google (Dev)' button
    - expect: User is redirected to the profile page
    - expect: URL changes to match pattern '/profile/*'
    - expect: Navigation bar shows 'Dev Google User' instead of 'Sign In'

#### 1.2. User menu navigation works correctly

**File:** `tests/authentication/user-menu.spec.ts`

**Steps:**
  1. Click on the 'Dev Google User' text in the navigation bar
    - expect: A dropdown menu appears
    - expect: Menu contains options: 'Profile', 'Find Films To Review', 'Film Graph', and 'Sign Out'
  2. Click on 'Profile' menu item
    - expect: User is navigated to their profile page
  3. Open the user menu again and click 'Find Films To Review'
    - expect: User is navigated to the recommendations page
  4. Open the user menu again and click 'Film Graph'
    - expect: User is navigated to the film graph page

#### 1.3. User can sign out

**File:** `tests/authentication/sign-out.spec.ts`

**Steps:**
  1. Click on the 'Dev Google User' text in the navigation bar
    - expect: Dropdown menu appears with 'Sign Out' option
  2. Click on 'Sign Out' menu item
    - expect: User is signed out
    - expect: Navigation bar shows 'Sign In' button instead of user name
    - expect: User is redirected to home page or login page

### 2. Film Search

**Seed:** `seed.spec.ts`

#### 2.1. User can search for films

**File:** `tests/film-search/search-films.spec.ts`

**Steps:**
  1. Navigate to the Search page by clicking 'Search' in the navigation bar
    - expect: URL changes to '/search'
    - expect: Page displays heading 'Search for Films'
    - expect: Search textbox is visible with placeholder 'Search for a film by title...'
    - expect: Empty state message shows 'Start Your Search' and 'Type in the search box above to find films'
  2. Type 'The Matrix' into the search textbox
    - expect: Search results appear automatically as user types
    - expect: Results heading shows 'Search Results (20)' or similar count
    - expect: Film cards display with poster images, titles, release dates, and descriptions
  3. Verify the first result is 'The Matrix (1999)'
    - expect: First film card shows title 'The Matrix'
    - expect: Release date shows '1999-03-31'
    - expect: Description is visible
  4. Clear the search and type 'Inception'
    - expect: Previous results are cleared
    - expect: New results appear for Inception
    - expect: Results heading shows updated count
    - expect: 'Inception (2010)' appears in the results

#### 2.2. Search handles no results gracefully

**File:** `tests/film-search/no-results.spec.ts`

**Steps:**
  1. Navigate to the Search page
    - expect: Search page loads successfully
  2. Type a nonsensical search term like 'xyzabc123notafilm'
    - expect: No results or appropriate message is displayed
    - expect: Application does not crash or show errors

#### 2.3. Search works with special characters

**File:** `tests/film-search/special-characters.spec.ts`

**Steps:**
  1. Navigate to the Search page
    - expect: Search page loads successfully
  2. Search for a film with special characters in the title (e.g., '12 Monkeys' or 'Amélie')
    - expect: Search returns relevant results
    - expect: Special characters are handled correctly

### 3. Film Review

**Seed:** `seed.spec.ts`

#### 3.1. User can submit a film review with rating

**File:** `tests/film-review/submit-review.spec.ts`

**Steps:**
  1. Navigate to Search page and search for 'The Matrix'
    - expect: Search results appear with The Matrix (1999)
  2. Click on The Matrix (1999) film card
    - expect: User is navigated to the review page for The Matrix
    - expect: URL changes to '/review/*'
    - expect: Page displays 'Review Film' heading
    - expect: Film poster, title, release date, and description are visible
    - expect: 5 star rating buttons are displayed
    - expect: Rating text shows '(Not rated)'
    - expect: Review textbox is visible with placeholder 'Share your thoughts about this film...'
    - expect: 'Submit Review' button is visible but disabled
  3. Click on the 5th star button to select 5-star rating
    - expect: Star button becomes active/highlighted
    - expect: Rating text updates to show '(5)'
  4. Type a review in the 'Your Thoughts' textbox: 'A groundbreaking sci-fi masterpiece that redefined action cinema.'
    - expect: Text appears in the textbox
    - expect: 'Submit Review' button becomes enabled
  5. Click the 'Submit Review' button
    - expect: Review is submitted successfully
    - expect: Success message appears: 'Review submitted successfully'
    - expect: Film's ELO rating section appears showing initial rating (e.g., 1,100)
    - expect: Rating comparison section appears
    - expect: Page heading changes to 'Update Review'
    - expect: Since this is the first film, comparison section shows 'No more films available for comparison'

#### 3.2. User can change star rating before submission

**File:** `tests/film-review/change-rating.spec.ts`

**Steps:**
  1. Navigate to review page for any film
    - expect: Review page loads with star rating buttons
  2. Click on the 3rd star to select 3-star rating
    - expect: Rating shows '(3)'
  3. Click on the 5th star to change rating to 5 stars
    - expect: Rating updates to show '(5)'
  4. Add review text and submit
    - expect: Review is submitted with 5-star rating

#### 3.3. Submit button is disabled without rating

**File:** `tests/film-review/validation-no-rating.spec.ts`

**Steps:**
  1. Navigate to review page for any film
    - expect: Review page loads
  2. Type review text without selecting a star rating
    - expect: 'Submit Review' button remains disabled
  3. Attempt to click the disabled submit button
    - expect: Nothing happens, form is not submitted

#### 3.4. Submit button is disabled without review text

**File:** `tests/film-review/validation-no-text.spec.ts`

**Steps:**
  1. Navigate to review page for any film
    - expect: Review page loads
  2. Select a star rating without typing review text
    - expect: 'Submit Review' button remains disabled
  3. Attempt to click the disabled submit button
    - expect: Nothing happens, form is not submitted

#### 3.5. Back to Search button works

**File:** `tests/film-review/back-to-search.spec.ts`

**Steps:**
  1. Navigate to any film review page
    - expect: Review page loads with 'Back to Search' button
  2. Click the 'Back to Search' button
    - expect: User is navigated back to the search page
    - expect: URL changes to '/search'
    - expect: Previous search term is cleared

### 4. Film Comparison

**Seed:** `seed.spec.ts`

#### 4.1. User can compare films using bulk mode

**File:** `tests/film-comparison/bulk-comparison.spec.ts`

**Steps:**
  1. Review and submit The Matrix with 5 stars
    - expect: Review submitted, comparison section shows no films available
  2. Navigate to search, find and review Inception with 4 stars
    - expect: Review submitted successfully
    - expect: Comparison section appears showing 'Rate Your Films' heading
    - expect: Bulk Mode checkbox is checked by default
    - expect: Progress indicator shows '1 / 50 films loaded'
    - expect: Target film shows 'Comparing: Inception (2010-07-15)'
    - expect: The Matrix film card appears with poster, title, and release date
    - expect: Three comparison buttons are displayed: 'Better', 'Same', 'Worse'
    - expect: 'Submit 0 Comparisons' button is disabled
  3. Click the 'Worse' button (meaning Inception is worse than The Matrix)
    - expect: 'Worse' button becomes active/highlighted
    - expect: 'Submit' button text updates to 'Submit 1 Comparisons' and becomes enabled
  4. Click the 'Submit 1 Comparisons' button
    - expect: Comparison is submitted successfully
    - expect: User is redirected to profile page
    - expect: Film Rankings section shows updated ELO ratings
    - expect: The Matrix has higher ELO rating than Inception (e.g., The Matrix: 1,120, Inception: 1,080)

#### 4.2. Bulk Mode checkbox toggles between modes

**File:** `tests/film-comparison/toggle-bulk-mode.spec.ts`

**Steps:**
  1. Navigate to a review page with comparison section visible (after reviewing second film)
    - expect: Comparison section shows with 'Bulk Mode' checkbox checked
    - expect: Multiple comparison buttons visible
  2. Click the 'Bulk Mode' checkbox to uncheck it
    - expect: Mode switches to sequential comparison mode
    - expect: UI adjusts accordingly for single film comparison
  3. Click the 'Bulk Mode' checkbox again to check it
    - expect: Mode switches back to bulk comparison mode

#### 4.3. User can select multiple comparisons before submitting

**File:** `tests/film-comparison/multiple-selections.spec.ts`

**Steps:**
  1. Review three films: The Matrix (5 stars), Inception (4 stars), The Dark Knight (5 stars)
    - expect: All three reviews submitted successfully
  2. Navigate to review page for a fourth film and submit review
    - expect: Comparison section appears
    - expect: Multiple films available for comparison
    - expect: Progress shows films loaded (e.g., '3 / 50 films loaded')
  3. Select 'Better' for first film, 'Worse' for second film, 'Same' for third film
    - expect: Each selection is highlighted
    - expect: Submit button updates to show 'Submit 3 Comparisons'
  4. Click submit button
    - expect: All comparisons are submitted at once
    - expect: User is redirected to profile page
    - expect: ELO ratings are updated for all compared films

#### 4.4. User can submit partial comparisons

**File:** `tests/film-comparison/partial-submission.spec.ts`

**Steps:**
  1. Navigate to comparison section with multiple films available
    - expect: Multiple film cards shown for comparison
  2. Select comparison choice for only 2 out of 10 available films
    - expect: Submit button shows 'Submit 2 Comparisons' and is enabled
  3. Click submit button
    - expect: Only the 2 selected comparisons are submitted
    - expect: Submission succeeds
    - expect: Profile page shows updated rankings

### 5. User Profile

**Seed:** `seed.spec.ts`

#### 5.1. Profile page displays user information

**File:** `tests/user-profile/view-profile.spec.ts`

**Steps:**
  1. After signing in, navigate to profile page via user menu or direct URL
    - expect: Profile page loads with URL pattern '/profile/*'
    - expect: User avatar/image is displayed
    - expect: User name 'Dev Google User' is shown as heading
    - expect: Username '@devgoogleuser' is displayed
    - expect: Member since date shows 'Member since January 2026'

#### 5.2. Recently Reviewed Films section displays correctly

**File:** `tests/user-profile/recently-reviewed-films.spec.ts`

**Steps:**
  1. View profile page before reviewing any films
    - expect: 'Recently Reviewed Films' section is visible
    - expect: Empty state shows 'No Reviews Yet'
    - expect: Message says 'Start reviewing films to see them appear here'
  2. Review a film (e.g., Inception with 4 stars)
    - expect: Review submitted successfully
  3. Navigate back to profile page
    - expect: 'Recently Reviewed Films' section shows the film
    - expect: Film card displays poster image, title 'Inception', review date 'Jan 30, 2026'
    - expect: 4 filled stars are shown
    - expect: Review text is displayed
    - expect: Navigation arrows appear if multiple films are reviewed

#### 5.3. Films to Review section displays correctly

**File:** `tests/user-profile/films-to-review.spec.ts`

**Steps:**
  1. View profile page with all reviewed films having comparisons
    - expect: 'Films to Review' section is visible
    - expect: Shows 'No Films Waiting'
    - expect: Message says 'All your seen films have been reviewed!'

#### 5.4. Film Rankings table displays and sorts correctly

**File:** `tests/user-profile/film-rankings.spec.ts`

**Steps:**
  1. Review and compare two films so ELO rankings are generated
    - expect: Films have ELO ratings
  2. Navigate to profile page and view 'Film Rankings' section
    - expect: Rankings table is displayed
    - expect: Column headers show 'Film' and 'Elo Rating'
    - expect: Both columns have sort icons
    - expect: Search box is visible with placeholder 'Search films...'
    - expect: Films are listed in order: '1. The Matrix 1,120' and '2. Inception 1,080'
    - expect: Pagination controls are visible at bottom
  3. Click on 'Film' column header to sort
    - expect: Films are sorted alphabetically by title
  4. Click on 'Elo Rating' column header to sort
    - expect: Films are sorted by ELO rating

#### 5.5. Film Rankings search functionality works

**File:** `tests/user-profile/rankings-search.spec.ts`

**Steps:**
  1. After having multiple films ranked, navigate to profile page
    - expect: Film Rankings table shows all films
  2. Type 'Matrix' in the search box
    - expect: Table filters to show only films with 'Matrix' in the title
    - expect: Other films are hidden
  3. Clear search box
    - expect: All films are displayed again

#### 5.6. Films Needing More Comparisons table displays correctly

**File:** `tests/user-profile/films-needing-comparisons.spec.ts`

**Steps:**
  1. After reviewing and comparing two films, view profile page
    - expect: 'Films Needing More Comparisons' section is visible
    - expect: Table shows columns: 'Film', 'Number of Comparisons', 'Release Year'
    - expect: Both films appear with comparison count of 1
    - expect: The Matrix shows year 1999
    - expect: Inception shows year 2010
    - expect: Pagination controls are available

#### 5.7. Profile page pagination works

**File:** `tests/user-profile/pagination.spec.ts`

**Steps:**
  1. Review more than 10 films to exceed default pagination
    - expect: Profile page loads with multiple films
  2. Scroll to Film Rankings pagination controls
    - expect: Pagination shows page numbers
    - expect: 'Next Page' button is enabled
    - expect: Rows per page dropdown shows '10'
  3. Click 'Next Page' button
    - expect: Table shows next set of films
    - expect: 'Previous Page' button becomes enabled
  4. Change 'Rows per page' to a different value
    - expect: Table updates to show selected number of rows per page

### 6. Film Recommendations

**Seed:** `seed.spec.ts`

#### 6.1. User can access recommendations page and see instructions

**File:** `tests/film-recommendations/view-recommendations-page.spec.ts`

**Steps:**
  1. Navigate to recommendations page via user menu 'Find Films To Review'
    - expect: URL changes to '/recommendations/*'
    - expect: Page heading shows 'Find Films to Review'
    - expect: Instructions heading: 'Select 3 Films You've Seen'
    - expect: Description text explains the purpose: 'These films will help us understand your taste...'
    - expect: 'Selected Films (0/3)' section is visible
    - expect: Empty state shows 'Search and add 3 films below to get started'
    - expect: 'Search for Films' section with search textbox is visible
    - expect: 'Start Generating Recommendations' button is disabled

#### 6.2. User can search and add films to recommendation seed list

**File:** `tests/film-recommendations/add-films.spec.ts`

**Steps:**
  1. Navigate to recommendations page
    - expect: Page loads with search functionality
  2. Type 'The Dark Knight' in the search textbox
    - expect: Search results appear below the textbox
    - expect: The Dark Knight (2008) appears in results with poster and release date
  3. Click on The Dark Knight (2008) film card
    - expect: Film is added to 'Selected Films' section
    - expect: Count updates to 'Selected Films (1/3)'
    - expect: Film card shows in selected section with remove button
    - expect: Search box is cleared
  4. Search for and add two more films (e.g., 'Pulp Fiction' and 'The Shawshank Redemption')
    - expect: Each film is added successfully
    - expect: Count updates: (2/3), then (3/3)
    - expect: 'Start Generating Recommendations' button becomes enabled

#### 6.3. User can remove films from selection

**File:** `tests/film-recommendations/remove-films.spec.ts`

**Steps:**
  1. Add a film to the selected films list
    - expect: Film appears in 'Selected Films' section with count (1/3)
  2. Click the remove button (X or trash icon) on the selected film
    - expect: Film is removed from selection
    - expect: Count updates back to 'Selected Films (0/3)'
    - expect: 'Start Generating Recommendations' button becomes disabled

#### 6.4. Cannot add the same film twice

**File:** `tests/film-recommendations/duplicate-prevention.spec.ts`

**Steps:**
  1. Add a film to the selected list (e.g., The Dark Knight)
    - expect: Film is added, count shows (1/3)
  2. Search for the same film again and attempt to add it
    - expect: Film either doesn't appear in search results, or clicking it has no effect
    - expect: Count remains at (1/3)
    - expect: Duplicate is not added

#### 6.5. Generate recommendations button requires 3 films

**File:** `tests/film-recommendations/button-validation.spec.ts`

**Steps:**
  1. Navigate to recommendations page
    - expect: 'Start Generating Recommendations' button is disabled
  2. Add only 1 film
    - expect: Button remains disabled
  3. Add second film
    - expect: Button remains disabled
  4. Add third film
    - expect: Button becomes enabled
  5. Remove one film to go back to 2 films
    - expect: Button becomes disabled again

### 7. Navigation

**Seed:** `seed.spec.ts`

#### 7.1. Main navigation links work correctly

**File:** `tests/navigation/main-nav-links.spec.ts`

**Steps:**
  1. From any page, click the '_cinema.log()' logo link
    - expect: User is navigated to home page
    - expect: URL changes to '/home'
    - expect: Logo link shows as active in navigation
  2. Click the 'Search' link in navigation
    - expect: User is navigated to search page
    - expect: URL changes to '/search'
    - expect: Search link shows as active

#### 7.2. Browser back and forward buttons work

**File:** `tests/navigation/browser-navigation.spec.ts`

**Steps:**
  1. Navigate from Home -> Search -> Review page
    - expect: Each navigation works correctly
  2. Click browser back button
    - expect: User returns to Search page
    - expect: Page state is preserved
  3. Click browser back button again
    - expect: User returns to Home page
  4. Click browser forward button
    - expect: User goes forward to Search page

#### 7.3. Active navigation state is correct

**File:** `tests/navigation/active-state.spec.ts`

**Steps:**
  1. Navigate to home page
    - expect: '_cinema.log()' link appears as active in navigation
  2. Navigate to search page
    - expect: 'Search' link appears as active
    - expect: Home link no longer shows as active

### 8. Error Handling

**Seed:** `seed.spec.ts`

#### 8.1. Application handles network errors gracefully

**File:** `tests/error-handling/network-errors.spec.ts`

**Steps:**
  1. Simulate network failure while searching for films (stop backend server or block network)
    - expect: Application shows appropriate error message
    - expect: No unhandled errors in console
    - expect: Application remains functional

#### 8.2. Invalid film ID shows error page

**File:** `tests/error-handling/invalid-film-id.spec.ts`

**Steps:**
  1. Navigate directly to a review URL with invalid/non-existent film ID: '/review/invalid-id-123'
    - expect: Application handles gracefully with error message or redirect
    - expect: User is not left on a broken page

#### 8.3. Session expiration is handled

**File:** `tests/error-handling/session-expiration.spec.ts`

**Steps:**
  1. Sign in and navigate to authenticated page
    - expect: User is authenticated and can access content
  2. Manually clear authentication cookies or wait for session to expire
    - expect: Next action requiring authentication shows appropriate message or redirects to login
