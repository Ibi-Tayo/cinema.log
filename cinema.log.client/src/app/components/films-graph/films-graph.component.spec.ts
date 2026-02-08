import { ComponentFixture, TestBed } from '@angular/core/testing';
import { describe, beforeEach, it, expect } from 'vitest';

import { FilmsGraphComponent } from './films-graph.component';

describe('FilmsGraphComponent', () => {
  let component: FilmsGraphComponent;
  let fixture: ComponentFixture<FilmsGraphComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FilmsGraphComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FilmsGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
