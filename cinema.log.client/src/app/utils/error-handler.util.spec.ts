import { HttpErrorResponse } from '@angular/common/http';
import {
  isExpectedError,
  handleHttpError,
  handleExpectedError,
} from './error-handler.util';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { firstValueFrom } from 'rxjs';

describe('Error Handler Utility', () => {
  let consoleErrorSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    consoleErrorSpy.mockRestore();
  });
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
    it('should log unexpected errors to console', async () => {
      const error = new HttpErrorResponse({ status: 500 });
      const handler = handleHttpError('fetching data', 'Failed to fetch data');

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).toHaveBeenCalledWith(
          'Error fetching data:',
          error,
        );
        expect(err.message).toBe('Failed to fetch data');
      }
    });

    it('should not log expected errors (404) to console', async () => {
      const error = new HttpErrorResponse({ status: 404 });
      const handler = handleHttpError('fetching data', 'Failed to fetch data');

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).not.toHaveBeenCalled();
        expect(err.message).toBe('Failed to fetch data');
      }
    });

    it('should not log expected errors (401) to console', async () => {
      const error = new HttpErrorResponse({ status: 401 });
      const handler = handleHttpError(
        'authenticating',
        'Authentication failed',
      );

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).not.toHaveBeenCalled();
        expect(err.message).toBe('Authentication failed');
      }
    });

    it('should not log expected errors (403) to console', async () => {
      const error = new HttpErrorResponse({ status: 403 });
      const handler = handleHttpError('accessing resource', 'Access denied');

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).not.toHaveBeenCalled();
        expect(err.message).toBe('Access denied');
      }
    });

    it('should not log when silent mode is enabled', async () => {
      const error = new HttpErrorResponse({ status: 500 });
      const handler = handleHttpError(
        'fetching data',
        'Failed to fetch data',
        true,
      );

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).not.toHaveBeenCalled();
        expect(err.message).toBe('Failed to fetch data');
      }
    });
  });

  describe('handleExpectedError', () => {
    it('should not log any errors to console', async () => {
      const error = new HttpErrorResponse({ status: 500 });
      const handler = handleExpectedError('Resource not found');

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).not.toHaveBeenCalled();
        expect(err.message).toBe('Resource not found');
      }
    });

    it('should handle 404 errors silently', async () => {
      const error = new HttpErrorResponse({ status: 404 });
      const handler = handleExpectedError('Resource not found');

      try {
        await firstValueFrom(handler(error));
        expect.fail('Should have thrown an error');
      } catch (err: any) {
        expect(consoleErrorSpy).not.toHaveBeenCalled();
        expect(err.message).toBe('Resource not found');
      }
    });
  });
});
