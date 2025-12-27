/**
 * TMDB Image Utility
 * Helper functions for constructing TMDB image URLs with various sizes
 */

export type TMDBPosterSize =
  | 'w92'
  | 'w154'
  | 'w185'
  | 'w342'
  | 'w500'
  | 'w780'
  | 'original';
export type TMDBBackdropSize = 'w300' | 'w780' | 'w1280' | 'original';

const TMDB_IMAGE_BASE_URL = 'https://image.tmdb.org/t/p/';

/**
 * Constructs a full TMDB poster URL from a poster path
 * @param posterPath - The poster path from TMDB API (e.g., "/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg")
 * @param size - The desired image size (default: 'w500')
 * @returns Full URL to the poster image, or empty string if posterPath is invalid
 */
export function getTMDBPosterUrl(
  posterPath: string | null | undefined,
  size: TMDBPosterSize = 'w500'
): string {
  if (!posterPath || posterPath.trim() === '') {
    return '';
  }

  // Remove leading slash if present (TMDB paths include it)
  const cleanPath = posterPath.startsWith('/') ? posterPath : `/${posterPath}`;

  return `${TMDB_IMAGE_BASE_URL}${size}${cleanPath}`;
}

/**
 * Constructs a full TMDB backdrop URL from a backdrop path
 * @param backdropPath - The backdrop path from TMDB API
 * @param size - The desired image size (default: 'w1280')
 * @returns Full URL to the backdrop image, or empty string if backdropPath is invalid
 */
export function getTMDBBackdropUrl(
  backdropPath: string | null | undefined,
  size: TMDBBackdropSize = 'w1280'
): string {
  if (!backdropPath || backdropPath.trim() === '') {
    return '';
  }

  // Remove leading slash if present (TMDB paths include it)
  const cleanPath = backdropPath.startsWith('/')
    ? backdropPath
    : `/${backdropPath}`;

  return `${TMDB_IMAGE_BASE_URL}${size}${cleanPath}`;
}
