import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LoadingStateComponent } from './loading-state.component';

describe('LoadingStateComponent', () => {
  let component: LoadingStateComponent;
  let fixture: ComponentFixture<LoadingStateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LoadingStateComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(LoadingStateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display default message when no input provided', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const messageElement = compiled.querySelector('p');
    expect(messageElement?.textContent).toContain('Loading...');
  });

  it('should display custom message when input provided', () => {
    fixture = TestBed.createComponent(LoadingStateComponent);
    fixture.componentRef.setInput('message', 'Loading film details...');
    fixture.detectChanges();

    const compiled = fixture.nativeElement as HTMLElement;
    const messageElement = compiled.querySelector('p');
    expect(messageElement?.textContent).toContain('Loading film details...');
  });

  it('should render lord-icon element', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const lordIcon = compiled.querySelector('lord-icon');
    expect(lordIcon).toBeTruthy();
    expect(lordIcon?.getAttribute('trigger')).toBe('loop');
  });

  it('should have loading-state class on container', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    const container = compiled.querySelector('.loading-state');
    expect(container).toBeTruthy();
  });
});
