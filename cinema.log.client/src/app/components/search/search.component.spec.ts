import type { MockedObject } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { of } from 'rxjs';
import { SearchComponent } from './search.component';
import { FilmService } from '../../services/film.service';

describe('SearchComponent', () => {
  let component: SearchComponent;
  let fixture: ComponentFixture<SearchComponent>;
  let mockFilmService: Partial<MockedObject<FilmService>>;
  let mockRouter: Partial<MockedObject<Router>>;

  beforeEach(async () => {
    mockFilmService = {
      searchFilms: vi.fn().mockName('FilmService.searchFilms'),
    };
    mockRouter = {
      navigate: vi.fn().mockName('Router.navigate'),
    };

    await TestBed.configureTestingModule({
      imports: [SearchComponent],
      providers: [
        { provide: FilmService, useValue: mockFilmService },
        { provide: Router, useValue: mockRouter },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(SearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
