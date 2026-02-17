import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { describe, beforeEach, it, expect } from 'vitest';

import { HomeComponent } from './home.component';
import { AuthService } from '../../services/auth.service';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let authService: AuthService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HomeComponent],
      providers: [provideHttpClient(), provideHttpClientTesting()],
    }).compileComponents();

    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    authService = TestBed.inject(AuthService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show CTA section when user is not logged in', () => {
    authService.currentUser.set(null);
    fixture.detectChanges();

    const compiled = fixture.nativeElement;
    const ctaSection = compiled.querySelector('[data-testid="home-cta-section"]');
    expect(ctaSection).toBeTruthy();
  });

  it('should hide CTA section when user is logged in', () => {
    authService.currentUser.set({
      id: '1',
      username: 'testuser',
      name: 'Test User',
      profilePicUrl: '',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    });
    fixture.detectChanges();

    const compiled = fixture.nativeElement;
    const ctaSection = compiled.querySelector('[data-testid="home-cta-section"]');
    expect(ctaSection).toBeNull();
  });
});
