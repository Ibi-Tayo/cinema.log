import { getTMDBPosterUrl, getTMDBBackdropUrl } from './tmdb-image.util';

describe('TMDB Image Utility', () => {
  describe('getTMDBPosterUrl', () => {
    it('should construct full URL with default size (w500)', () => {
      const posterPath = '/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg';
      const result = getTMDBPosterUrl(posterPath);
      expect(result).toBe(
        'https://image.tmdb.org/t/p/w500/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg'
      );
    });

    it('should construct full URL with custom size', () => {
      const posterPath = '/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg';
      const result = getTMDBPosterUrl(posterPath, 'w342');
      expect(result).toBe(
        'https://image.tmdb.org/t/p/w342/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg'
      );
    });

    it('should handle path without leading slash', () => {
      const posterPath = 'bOGkgRGdhrBYJSLpXaxhXVstddV.jpg';
      const result = getTMDBPosterUrl(posterPath, 'w185');
      expect(result).toBe(
        'https://image.tmdb.org/t/p/w185/bOGkgRGdhrBYJSLpXaxhXVstddV.jpg'
      );
    });

    it('should return empty string for null posterPath', () => {
      const result = getTMDBPosterUrl(null);
      expect(result).toBe('');
    });

    it('should return empty string for undefined posterPath', () => {
      const result = getTMDBPosterUrl(undefined);
      expect(result).toBe('');
    });

    it('should return empty string for empty posterPath', () => {
      const result = getTMDBPosterUrl('');
      expect(result).toBe('');
    });

    it('should return empty string for whitespace-only posterPath', () => {
      const result = getTMDBPosterUrl('   ');
      expect(result).toBe('');
    });

    it('should work with all poster sizes', () => {
      const posterPath = '/test.jpg';

      expect(getTMDBPosterUrl(posterPath, 'w92')).toBe(
        'https://image.tmdb.org/t/p/w92/test.jpg'
      );
      expect(getTMDBPosterUrl(posterPath, 'w154')).toBe(
        'https://image.tmdb.org/t/p/w154/test.jpg'
      );
      expect(getTMDBPosterUrl(posterPath, 'w185')).toBe(
        'https://image.tmdb.org/t/p/w185/test.jpg'
      );
      expect(getTMDBPosterUrl(posterPath, 'w342')).toBe(
        'https://image.tmdb.org/t/p/w342/test.jpg'
      );
      expect(getTMDBPosterUrl(posterPath, 'w500')).toBe(
        'https://image.tmdb.org/t/p/w500/test.jpg'
      );
      expect(getTMDBPosterUrl(posterPath, 'w780')).toBe(
        'https://image.tmdb.org/t/p/w780/test.jpg'
      );
      expect(getTMDBPosterUrl(posterPath, 'original')).toBe(
        'https://image.tmdb.org/t/p/original/test.jpg'
      );
    });
  });

  describe('getTMDBBackdropUrl', () => {
    it('should construct full URL with default size (w1280)', () => {
      const backdropPath = '/backdrop.jpg';
      const result = getTMDBBackdropUrl(backdropPath);
      expect(result).toBe('https://image.tmdb.org/t/p/w1280/backdrop.jpg');
    });

    it('should construct full URL with custom size', () => {
      const backdropPath = '/backdrop.jpg';
      const result = getTMDBBackdropUrl(backdropPath, 'w780');
      expect(result).toBe('https://image.tmdb.org/t/p/w780/backdrop.jpg');
    });

    it('should return empty string for null backdropPath', () => {
      const result = getTMDBBackdropUrl(null);
      expect(result).toBe('');
    });

    it('should work with all backdrop sizes', () => {
      const backdropPath = '/backdrop.jpg';

      expect(getTMDBBackdropUrl(backdropPath, 'w300')).toBe(
        'https://image.tmdb.org/t/p/w300/backdrop.jpg'
      );
      expect(getTMDBBackdropUrl(backdropPath, 'w780')).toBe(
        'https://image.tmdb.org/t/p/w780/backdrop.jpg'
      );
      expect(getTMDBBackdropUrl(backdropPath, 'w1280')).toBe(
        'https://image.tmdb.org/t/p/w1280/backdrop.jpg'
      );
      expect(getTMDBBackdropUrl(backdropPath, 'original')).toBe(
        'https://image.tmdb.org/t/p/original/backdrop.jpg'
      );
    });
  });
});
