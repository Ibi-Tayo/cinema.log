import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { UserService, User } from './user.service';
import { environment } from '../../environments/environment';

describe('UserService', () => {
  let service: UserService;
  let httpMock: HttpTestingController;

  const mockUser: User = {
    id: '123e4567-e89b-12d3-a456-426614174000',
    githubId: 12345,
    name: 'Test User',
    username: 'testuser',
    profilePicUrl: 'https://example.com/avatar.png',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [UserService],
    });

    service = TestBed.inject(UserService);
    httpMock = TestBed.inject(HttpTestingController);
    spyOn(console, 'error');
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should get user by id', () => {
    service
      .getUserById('123e4567-e89b-12d3-a456-426614174000')
      .subscribe((user) => {
        expect(user).toEqual(mockUser);
      });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/users/123e4567-e89b-12d3-a456-426614174000`
    );
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockUser);
  });

  it('should handle error when getting user by id', () => {
    service.getUserById('invalid-id').subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to fetch user');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users/invalid-id`);
    req.error(new ProgressEvent('error'));
  });

  it('should get all users', () => {
    const mockUsers: User[] = [mockUser];

    service.getAllUsers().subscribe((users) => {
      expect(users).toEqual(mockUsers);
      expect(users.length).toBe(1);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users`);
    expect(req.request.method).toBe('GET');
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockUsers);
  });

  it('should create user', () => {
    const newUser: Partial<User> = {
      name: 'New User',
      username: 'newuser',
    };

    service.createUser(newUser).subscribe((user) => {
      expect(user).toEqual(mockUser);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(newUser);
    expect(req.request.withCredentials).toBe(true);
    req.flush(mockUser);
  });

  it('should handle error when creating user', () => {
    const newUser: Partial<User> = {
      name: 'New User',
      username: 'newuser',
    };

    service.createUser(newUser).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to create user');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users`);
    req.error(new ProgressEvent('error'));
  });

  it('should update user', () => {
    const updatedUser: User = { ...mockUser, name: 'Updated User' };

    service.updateUser(updatedUser).subscribe((user) => {
      expect(user).toEqual(updatedUser);
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(updatedUser);
    expect(req.request.withCredentials).toBe(true);
    req.flush(updatedUser);
  });

  it('should handle error when updating user', () => {
    service.updateUser(mockUser).subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to update user');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users`);
    req.error(new ProgressEvent('error'));
  });

  it('should delete user', () => {
    service.deleteUser('123e4567-e89b-12d3-a456-426614174000').subscribe(() => {
      expect(true).toBe(true);
    });

    const req = httpMock.expectOne(
      `${environment.apiUrl}/users/123e4567-e89b-12d3-a456-426614174000`
    );
    expect(req.request.method).toBe('DELETE');
    expect(req.request.withCredentials).toBe(true);
    req.flush(null);
  });

  it('should handle error when deleting user', () => {
    service.deleteUser('invalid-id').subscribe({
      next: () => fail('should have failed'),
      error: (error) => {
        expect(error.message).toContain('Failed to delete user');
      },
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/users/invalid-id`);
    req.error(new ProgressEvent('error'));
  });
});
