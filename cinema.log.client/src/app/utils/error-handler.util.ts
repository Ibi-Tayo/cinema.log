import { HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';

/**
 * Determines if an HTTP error is expected in normal application flow.
 * Expected errors include:
 * - 404: Resource not found (e.g., rating or review doesn't exist yet)
 * - 401: Unauthorized (e.g., user not logged in)
 * - 403: Forbidden (e.g., user doesn't have permission)
 */
export function isExpectedError(error: any): boolean {
  if (error instanceof HttpErrorResponse) {
    return error.status === 404 || error.status === 401 || error.status === 403;
  }
  return false;
}

/**
 * Handles HTTP errors gracefully by:
 * - Logging unexpected errors to console
 * - Silently handling expected errors (404, 401, 403)
 * - Returning an observable error with a user-friendly message
 * 
 * @param context - Description of the operation for logging (e.g., 'fetching rating')
 * @param errorMessage - User-friendly error message to return
 * @param silent - If true, suppresses all console logging (for expected 404s, 401s, etc.)
 */
export function handleHttpError(
  context: string,
  errorMessage: string,
  silent = false
): (error: any) => Observable<never> {
  return (error: any) => {
    // Only log to console if:
    // 1. Not in silent mode, AND
    // 2. It's not an expected error (404, 401, 403)
    if (!silent && !isExpectedError(error)) {
      console.error(`Error ${context}:`, error);
    }
    
    return throwError(() => new Error(errorMessage));
  };
}

/**
 * Handles expected errors that should not be logged to console.
 * Use this for operations where the absence of a resource is part of normal flow.
 * 
 * @param errorMessage - User-friendly error message to return
 */
export function handleExpectedError(errorMessage: string): (error: any) => Observable<never> {
  return handleHttpError('', errorMessage, true);
}
