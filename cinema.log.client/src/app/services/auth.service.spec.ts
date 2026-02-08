import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { AuthService } from './auth.service';
import { Router } from '@angular/router';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;
  let routerSpy: jasmine.SpyObj<Router>;

  beforeEach(() => {
    const routerSpyObj = jasmine.createSpyObj('Router', ['navigate']);

    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [AuthService, { provide: Router, useValue: routerSpyObj }],
    });

    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
    routerSpy = TestBed.inject(Router) as jasmine.SpyObj<Router>;
    spyOn(console, 'error');
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should logout successfully', () => {
    service.logout().subscribe({
      next: () => {
        expect(true).toBe(true);
      },
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/auth/logout`,
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush({});
  });

  it('should handle logout error', () => {
    service.logout().subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Logout failed');
      },
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/auth/logout`,
    );
    req.error(new ProgressEvent('error'));
  });

  it('should request refresh token successfully', () => {
    service.requestRefreshToken().subscribe({
      next: () => {
        expect(true).toBe(true);
      },
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/auth/refresh-token`,
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush({});
  });

  it('should handle refresh token error', () => {
    service.requestRefreshToken().subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Authentication session expired');
      },
    });

    const req = httpMock.expectOne(
      `${import.meta.env.NG_APP_API_URL}/auth/refresh-token`,
    );
    req.error(new ProgressEvent('error'));
  });

  it('should parse cookies correctly', () => {
    // Mock document.cookie
    Object.defineProperty(document, 'cookie', {
      writable: true,
      value: 'test-cookie=test-value; another-cookie=another-value',
    });

    const cookieValue = service.getCookie('test-cookie');
    expect(cookieValue).toBe('test-value');
  });

  it('should return empty string for non-existent cookie', () => {
    Object.defineProperty(document, 'cookie', {
      writable: true,
      value: 'test-cookie=test-value',
    });

    const cookieValue = service.getCookie('non-existent');
    expect(cookieValue).toBe('');
  });
});
