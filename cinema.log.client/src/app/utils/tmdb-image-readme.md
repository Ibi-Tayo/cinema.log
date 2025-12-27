# TMDB Image Utility

This utility provides helper functions for constructing TMDB (The Movie Database) image URLs with configurable sizes.

## Overview

TMDB API returns image paths (e.g., `/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg`) that need to be combined with a base URL and size parameter to create the full image URL.

## Usage

### Basic Usage

```typescript
import { getTMDBPosterUrl } from "../../utils/tmdb-image.util";

// In your component
const posterUrl = getTMDBPosterUrl(film.posterUrl); // Uses default size 'w500'
```

### With Custom Size

```typescript
import { getTMDBPosterUrl, TMDBPosterSize } from "../../utils/tmdb-image.util";

// Use smaller size for thumbnails
const thumbnailUrl = getTMDBPosterUrl(film.posterUrl, "w185");

// Use larger size for detail views
const detailUrl = getTMDBPosterUrl(film.posterUrl, "w780");

// Use original size (largest available)
const originalUrl = getTMDBPosterUrl(film.posterUrl, "original");
```

### In Templates

```html
<!-- Using default size -->
<img [src]="getPosterUrl(film.posterUrl)" [alt]="film.title" />

<!-- Using custom size -->
<img [src]="getPosterUrl(film.posterUrl, 'w342')" [alt]="film.title" />
```

## Available Poster Sizes

The following sizes are available for posters:

- `w92` - 92px wide (very small thumbnails)
- `w154` - 154px wide (small thumbnails)
- `w185` - 185px wide (medium thumbnails)
- `w342` - 342px wide (large thumbnails)
- `w500` - 500px wide (default, medium detail)
- `w780` - 780px wide (high detail)
- `original` - Original resolution (largest available)

## Size Recommendations

- **Search Results / Grid View**: `w342` - Good balance for grid layouts
- **Detail View / Single Film**: `w500` or `w780` - Higher quality for prominent display
- **Small Thumbnails**: `w185` - For lists or carousels
- **Full Screen / Hero Images**: `original` - Maximum quality

## Backdrop Images

For backdrop/fanart images, use `getTMDBBackdropUrl`:

```typescript
import { getTMDBBackdropUrl } from "../../utils/tmdb-image.util";

const backdropUrl = getTMDBBackdropUrl(film.backdropPath, "w1280");
```

Available backdrop sizes: `w300`, `w780`, `w1280`, `original`

## Implementation in Components

Each component has a `getPosterUrl` helper method:

```typescript
/**
 * Gets the TMDB poster URL for a film with specified size
 */
getPosterUrl(posterPath: string | null | undefined, size: TMDBPosterSize = 'w500'): string {
  return getTMDBPosterUrl(posterPath, size);
}
```

This allows templates to easily construct URLs while keeping the utility logic centralized.

## Error Handling

The utility handles edge cases gracefully:

- Returns empty string for `null`, `undefined`, or empty poster paths
- Automatically adds leading slash if missing from the path
- No errors thrown for invalid input

## Performance Considerations

- Smaller sizes load faster and use less bandwidth
- Use appropriate sizes for the display context to optimize performance
- Consider using `w342` or `w500` for most use cases
- Only use `original` when maximum quality is required
