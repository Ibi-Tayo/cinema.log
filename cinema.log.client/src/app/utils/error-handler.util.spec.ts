import { HttpErrorResponse } from '@angular/common/http';
import { isExpectedError, handleHttpError, handleExpectedError } from './error-handler.util';

describe('Error Handler Utility', () => {
  describe('isExpectedError', () => {
    it('should return true for 404 errors', () => {
      const error = new HttpErrorResponse({ status: 404 });
      expect(isExpectedError(error)).toBe(true);
    });

    it('should return true for 401 errors', () => {
      const error = new HttpErrorResponse({ status: 401 });
      expect(isExpectedError(error)).toBe(true);
    });

    it('should return true for 403 errors', () => {
      const error = new HttpErrorResponse({ status: 403 });
      expect(isExpectedError(error)).toBe(true);
    });

    it('should return false for 500 errors', () => {
      const error = new HttpErrorResponse({ status: 500 });
      expect(isExpectedError(error)).toBe(false);
    });

    it('should return false for 400 errors', () => {
      const error = new HttpErrorResponse({ status: 400 });
      expect(isExpectedError(error)).toBe(false);
    });

    it('should return false for non-HTTP errors', () => {
      const error = new Error('Network error');
      expect(isExpectedError(error)).toBe(false);
    });
  });

  describe('handleHttpError', () => {
    it('should log unexpected errors to console', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 500 });
      const handler = handleHttpError('fetching data', 'Failed to fetch data');

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).toHaveBeenCalledWith('Error fetching data:', error);
          expect(err.message).toBe('Failed to fetch data');
          done();
        },
      });
    });

    it('should not log expected errors (404) to console', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 404 });
      const handler = handleHttpError('fetching data', 'Failed to fetch data');

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).not.toHaveBeenCalled();
          expect(err.message).toBe('Failed to fetch data');
          done();
        },
      });
    });

    it('should not log expected errors (401) to console', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 401 });
      const handler = handleHttpError('authenticating', 'Authentication failed');

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).not.toHaveBeenCalled();
          expect(err.message).toBe('Authentication failed');
          done();
        },
      });
    });

    it('should not log expected errors (403) to console', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 403 });
      const handler = handleHttpError('accessing resource', 'Access denied');

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).not.toHaveBeenCalled();
          expect(err.message).toBe('Access denied');
          done();
        },
      });
    });

    it('should not log when silent mode is enabled', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 500 });
      const handler = handleHttpError('fetching data', 'Failed to fetch data', true);

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).not.toHaveBeenCalled();
          expect(err.message).toBe('Failed to fetch data');
          done();
        },
      });
    });
  });

  describe('handleExpectedError', () => {
    it('should not log any errors to console', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 500 });
      const handler = handleExpectedError('Resource not found');

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).not.toHaveBeenCalled();
          expect(err.message).toBe('Resource not found');
          done();
        },
      });
    });

    it('should handle 404 errors silently', (done) => {
      spyOn(console, 'error');
      const error = new HttpErrorResponse({ status: 404 });
      const handler = handleExpectedError('Resource not found');

      handler(error).subscribe({
        error: (err) => {
          expect(console.error).not.toHaveBeenCalled();
          expect(err.message).toBe('Resource not found');
          done();
        },
      });
    });
  });
});
