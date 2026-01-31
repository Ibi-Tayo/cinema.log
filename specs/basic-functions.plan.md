# Cinema.log Basic Functions Test Plan

## Application Overview

Cinema.log is a full-stack film rating application that allows users to review films, compare them using an ELO rating system, and build personalized rankings. The application features GitHub OAuth authentication, film search via TMDB API, bulk comparison mode, and visual film graph representation. This test plan covers all core functionality including authentication, film search, review creation, comparison workflows, profile management, and navigation.

## Test Scenarios

### 1. Authentication

**Seed:** `tests/seed.spec.ts`

#### 1.1. User can sign in successfully

**File:** `tests/authentication/sign-in.spec.ts`

**Steps:**
  1. Navigate to the home page at http://localhost:4200/
    - expect: The page loads successfully
    - expect: The home page displays the heading 'Your personal hub for film review'
  2. Click the 'Sign In' button in the navigation bar
    - expect: The login page appears
    - expect: Authentication options are displayed
  3. Click the 'Sign in with Google (Dev)' button
    - expect: The user is redirected to their profile page
    - expect: The URL changes to /profile/{user-id}
    - expect: The navigation bar displays 'Dev Google User' instead of 'Sign In'

#### 1.2. User can sign out successfully

**File:** `tests/authentication/sign-out.spec.ts`

**Steps:**
  1. Click the 'Dev Google User' dropdown in the navigation bar
    - expect: A dropdown menu appears with options: Profile, Find Films To Review, Film Graph, Sign Out
  2. Click 'Sign Out' from the dropdown menu
    - expect: The user is logged out
    - expect: The user is redirected to the home page
    - expect: The navigation bar shows 'Sign In' button

#### 1.3. Unauthenticated user cannot access protected pages

**File:** `tests/authentication/protected-routes.spec.ts`

**Steps:**
  1. Navigate directly to a profile URL without being logged in (e.g., /profile/c55bda9d-5434-4df6-81c2-52a9e0d63c55)
    - expect: The user is redirected to the login page or home page
    - expect: Access to the profile page is denied
  2. Navigate directly to the search page without being logged in
    - expect: The user is redirected to the login page or home page
  3. Navigate directly to the recommendations page without being logged in
    - expect: The user is redirected to the login page or home page

### 2. Film Search

**Seed:** `tests/seed.spec.ts`

#### 2.1. User can search for films by title

**File:** `tests/film-search/search-by-title.spec.ts`

**Steps:**
  1. Click 'Search' in the navigation bar
    - expect: The search page loads
    - expect: The heading 'Search for Films' is displayed
    - expect: An empty search textbox with placeholder 'Search for a film by title...' is visible
    - expect: A message 'Start Your Search' is displayed
  2. Type 'The Dark Knight' in the search textbox
    - expect: The text appears in the search field
  3. Press Enter or wait for search results
    - expect: Search results appear showing films matching the query
    - expect: The heading 'Search Results (X)' shows the number of results
    - expect: The first result shows 'The Dark Knight' (2008) with poster image and description
    - expect: Related films like 'The Dark Knight Rises' also appear in results

#### 2.2. Search displays empty state for no results

**File:** `tests/film-search/no-results.spec.ts`

**Steps:**
  1. Navigate to the search page
    - expect: The search page loads successfully
  2. Type a nonsensical string that will return no results (e.g., 'xyzabc123nonexistentfilm')
    - expect: The search executes
  3. Wait for search completion
    - expect: A message indicating no results were found is displayed
    - expect: No film cards are shown

#### 2.3. User can click on a search result to review

**File:** `tests/film-search/select-film.spec.ts`

**Steps:**
  1. Search for 'Inception'
    - expect: Search results display including the film 'Inception'
  2. Click on the 'Inception' film card from search results
    - expect: The user is redirected to the review page
    - expect: The URL changes to /review/{film-id}
    - expect: The page displays 'Review Film' heading
    - expect: Film details are shown including title 'Inception', release date, poster, and description

### 3. Film Review and Rating

**Seed:** `tests/seed.spec.ts`

#### 3.1. User can submit a complete film review

**File:** `tests/review/submit-review.spec.ts`

**Steps:**
  1. Search for and select a film not yet reviewed (e.g., 'Pulp Fiction')
    - expect: The review page loads with the film details
  2. Click the 5-star rating button
    - expect: The star is highlighted/selected
    - expect: The text updates to show 'Click to rate (5)'
  3. Type a review in the 'Your Thoughts' textbox: 'A masterpiece of storytelling with unforgettable dialogue'
    - expect: The text appears in the textarea
    - expect: The 'Submit Review' button becomes enabled
  4. Click the 'Submit Review' button
    - expect: A success message 'Review submitted successfully' appears
    - expect: The page transitions to show the comparison section
    - expect: The heading changes to 'Update Review'
    - expect: Film ELO rating is displayed showing 1,100 'Based on 0 comparisons'

#### 3.2. User cannot submit review without rating

**File:** `tests/review/validation-no-rating.spec.ts`

**Steps:**
  1. Navigate to review page for a new film
    - expect: The review form is displayed
  2. Type review text without selecting a star rating
    - expect: The 'Submit Review' button remains disabled
  3. Attempt to click the disabled 'Submit Review' button
    - expect: Nothing happens
    - expect: The review is not submitted

#### 3.3. User can select different star ratings

**File:** `tests/review/star-rating-selection.spec.ts`

**Steps:**
  1. Navigate to a film review page
    - expect: Five star rating buttons are displayed
  2. Click the 3-star button
    - expect: The rating text shows 'Click to rate (3)'
  3. Click the 5-star button
    - expect: The rating text updates to show 'Click to rate (5)'
    - expect: The previous selection is replaced
  4. Click the 1-star button
    - expect: The rating text updates to show 'Click to rate (1)'

#### 3.4. User can navigate back to search from review page

**File:** `tests/review/back-to-search.spec.ts`

**Steps:**
  1. Search for and select a film to review
    - expect: The review page loads
    - expect: A 'Back to Search' button is visible
  2. Click the 'Back to Search' button
    - expect: The user is navigated back to the search page
    - expect: Previous search results are still visible or search is in initial state

### 4. Bulk Film Comparison

**Seed:** `tests/seed.spec.ts`

#### 4.1. User can perform bulk comparisons

**File:** `tests/comparison/bulk-comparison.spec.ts`

**Steps:**
  1. After submitting a review, verify bulk mode is enabled
    - expect: The 'Bulk Mode' checkbox is checked
    - expect: The comparison section displays 'Select your preference for each film below'
    - expect: Multiple film cards are shown with Better/Same/Worse buttons
    - expect: A counter shows '{X} / 50 films loaded'
    - expect: The target film is displayed at the top: 'Comparing: {Film Title} ({Date})'
  2. Click 'Better' on the first comparison film (e.g., Inception)
    - expect: The button is visually selected/highlighted
    - expect: The comparison counter increments
  3. Click 'Worse' on the second comparison film (e.g., The Matrix)
    - expect: The button is visually selected/highlighted
    - expect: The submit button text updates to 'Submit 2 Comparisons'
    - expect: The submit button becomes enabled
  4. Click the 'Submit {X} Comparisons' button
    - expect: The comparisons are saved
    - expect: A success message appears
    - expect: The ELO ratings are updated
    - expect: The comparison count increases

#### 4.2. User can toggle between bulk and sequential mode

**File:** `tests/comparison/mode-toggle.spec.ts`

**Steps:**
  1. On the comparison screen, verify bulk mode is active
    - expect: The 'Bulk Mode' checkbox is checked
    - expect: Multiple film cards are displayed simultaneously
  2. Click the 'Bulk Mode' checkbox to uncheck it
    - expect: The mode switches to sequential
    - expect: Only one film comparison is shown at a time
    - expect: Navigation buttons (previous/next) appear for moving between films
  3. Click the checkbox again to re-enable bulk mode
    - expect: The mode switches back to bulk
    - expect: Multiple film cards are displayed again

#### 4.3. Bulk mode allows partial submissions

**File:** `tests/comparison/partial-submission.spec.ts`

**Steps:**
  1. Navigate to comparison screen with multiple films loaded
    - expect: Multiple film cards are displayed
  2. Select comparison preference for only 1 out of 3+ available films
    - expect: The submit button shows 'Submit 1 Comparison'
    - expect: The button becomes enabled
  3. Click the submit button
    - expect: The single comparison is saved successfully
    - expect: The submitted film is removed from comparison list
    - expect: Remaining uncompared films are still available

#### 4.4. User cannot submit without any selections

**File:** `tests/comparison/no-selection.spec.ts`

**Steps:**
  1. Navigate to comparison screen
    - expect: Film comparison cards are displayed
    - expect: The submit button shows 'Submit 0 Comparisons' and is disabled
  2. Attempt to click the disabled submit button
    - expect: Nothing happens
    - expect: No comparisons are submitted

### 5. User Profile

**Seed:** `tests/seed.spec.ts`

#### 5.1. Profile displays user information correctly

**File:** `tests/profile/profile-information.spec.ts`

**Steps:**
  1. Click the user dropdown and select 'Profile'
    - expect: The profile page loads
    - expect: User avatar/image is displayed
    - expect: Username 'Dev Google User' is shown as a heading
    - expect: User handle '@devgoogleuser' is displayed
    - expect: Member since date shows 'Member since January 2026'

#### 5.2. Recently Reviewed Films section displays correct data

**File:** `tests/profile/recent-reviews.spec.ts`

**Steps:**
  1. Navigate to profile page
    - expect: The 'Recently Reviewed Films' section is visible
  2. Verify the section content
    - expect: Recently reviewed films are displayed with poster images
    - expect: Each film shows: title, review date, star rating (visual stars), and review text
    - expect: Navigation arrows appear if multiple reviews exist
    - expect: Films are ordered by review date (most recent first)

#### 5.3. Film Rankings table displays and sorts correctly

**File:** `tests/profile/film-rankings.spec.ts`

**Steps:**
  1. Navigate to profile page and locate the 'Film Rankings' section
    - expect: A table is displayed with columns: Film, ELO Rating
    - expect: Films are listed with their rank numbers (1, 2, 3, etc.)
    - expect: Each row shows film name and ELO rating
    - expect: Films are sorted by ELO rating in descending order
  2. Click the 'Film' column header
    - expect: The table re-sorts alphabetically by film name
    - expect: A sort indicator icon appears on the column header
  3. Click the 'ELO Rating' column header
    - expect: The table re-sorts by ELO rating
    - expect: A sort indicator icon appears
  4. Type a film name in the search textbox
    - expect: The table filters to show only films matching the search query

#### 5.4. Films Needing More Comparisons table functions correctly

**File:** `tests/profile/needs-comparison.spec.ts`

**Steps:**
  1. Navigate to profile page and locate 'Films Needing More Comparisons' section
    - expect: A table displays with columns: Film, Number of Comparisons, Release Year
    - expect: Films with fewer comparisons are shown
    - expect: Each row is clickable
  2. Click on a film in the table
    - expect: The user is navigated to that film's comparison page
  3. Click different column headers to sort
    - expect: The table re-sorts based on the selected column
    - expect: Sort indicators update accordingly

#### 5.5. Profile pagination works correctly

**File:** `tests/profile/pagination.spec.ts`

**Steps:**
  1. Navigate to profile with sufficient films to enable pagination (>10 films)
    - expect: Tables show pagination controls
    - expect: Rows per page dropdown shows '10'
  2. Click the 'Next Page' button
    - expect: The next set of results loads
    - expect: Page indicator updates
    - expect: 'Previous Page' button becomes enabled
  3. Change the rows per page dropdown to a different value
    - expect: The table updates to show the selected number of rows per page

### 6. Film Graph

**Seed:** `tests/seed.spec.ts`

#### 6.1. User can access film graph visualization

**File:** `tests/film-graph/access-graph.spec.ts`

**Steps:**
  1. Click the user dropdown in navigation
    - expect: Dropdown menu appears with 'Film Graph' option
  2. Click 'Film Graph' from the dropdown
    - expect: The film graph page loads
    - expect: The URL changes to /film-graph
    - expect: The heading 'Films You Have Seen' is displayed

#### 6.2. Film graph displays user's film network

**File:** `tests/film-graph/graph-visualization.spec.ts`

**Steps:**
  1. Navigate to the film graph page
    - expect: A graph visualization is rendered
    - expect: Film nodes represent rated/reviewed films
    - expect: Connections between films are visible (if comparisons exist)

### 7. Recommendations

**Seed:** `tests/seed.spec.ts`

#### 7.1. User can access recommendations page

**File:** `tests/recommendations/access-recommendations.spec.ts`

**Steps:**
  1. Click user dropdown and select 'Find Films To Review'
    - expect: The recommendations page loads
    - expect: The URL changes to /recommendations/{user-id}
    - expect: The heading 'Find Films to Review' is displayed

#### 7.2. User can select seed films for recommendations

**File:** `tests/recommendations/select-seed-films.spec.ts`

**Steps:**
  1. Navigate to recommendations page
    - expect: The heading 'Select 3 Films You've Seen' is displayed
    - expect: Instructions about selecting variety of genres are shown
    - expect: 'Selected Films (0/3)' counter is displayed
    - expect: A search textbox with placeholder 'Search for a film...' is visible
    - expect: 'Start Generating Recommendations' button is disabled
  2. Search for a film in the search box
    - expect: Search results appear below the search box
  3. Click on a film from search results to add it
    - expect: The film is added to selected films
    - expect: Counter updates to 'Selected Films (1/3)'
  4. Repeat to add two more films
    - expect: Counter updates to (2/3) then (3/3)
    - expect: The 'Start Generating Recommendations' button becomes enabled

#### 7.3. User cannot generate recommendations with fewer than 3 films

**File:** `tests/recommendations/validation.spec.ts`

**Steps:**
  1. Navigate to recommendations page
    - expect: Button is disabled showing it requires selections
  2. Add only 1 or 2 films to selections
    - expect: The 'Start Generating Recommendations' button remains disabled
  3. Add a third film
    - expect: The button becomes enabled

### 8. Navigation

**Seed:** `tests/seed.spec.ts`

#### 8.1. Logo navigation returns user to home

**File:** `tests/navigation/logo-home.spec.ts`

**Steps:**
  1. Navigate to any page in the application (e.g., search, profile)
    - expect: The user is on a page other than home
  2. Click the '_cinema.log()' logo in the navigation bar
    - expect: The user is navigated to the home page
    - expect: The URL changes to /home
    - expect: The home page heading 'Your personal hub for film review' is displayed

#### 8.2. Navigation bar persists across all pages

**File:** `tests/navigation/navbar-persistence.spec.ts`

**Steps:**
  1. Navigate to home page
    - expect: Navigation bar is visible at the top with logo, Search link, and user menu
  2. Navigate to search page
    - expect: Navigation bar remains visible with same elements
  3. Navigate to profile page
    - expect: Navigation bar remains visible with same elements
  4. Navigate to film graph
    - expect: Navigation bar remains visible with same elements

#### 8.3. Active navigation items are highlighted

**File:** `tests/navigation/active-state.spec.ts`

**Steps:**
  1. Navigate to the Search page
    - expect: The 'Search' link in navigation has an active/highlighted state
  2. Click the logo to go home
    - expect: The logo/home link has an active state
    - expect: The Search link is no longer highlighted

### 9. Home Page

**Seed:** `tests/seed.spec.ts`

#### 9.1. Home page displays correct content

**File:** `tests/home/home-content.spec.ts`

**Steps:**
  1. Navigate to the home page
    - expect: The main heading 'Your personal hub for film review' is displayed
    - expect: A subheading 'A smarter way to track and rate the films you love' is shown
    - expect: Three main content sections are visible: 'Review & Rate Films', 'ELO Rating System', and 'Ready to Start?'
  2. Scroll to view the 'ELO Rating System' section
    - expect: A 'How It Works' explanation with three numbered steps is displayed
    - expect: Step 1: 'Compare Films' with description
    - expect: Step 2: 'Dynamic Ratings' with description
    - expect: Step 3: 'Refined Over Time' with description
  3. Review the 'Review & Rate Films' section
    - expect: Five star icons are displayed
    - expect: Text explains the 1-5 star rating system
