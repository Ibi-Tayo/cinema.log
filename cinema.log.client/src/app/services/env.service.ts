import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

declare global {
  interface Window {
    __env?: {
      apiUrl?: string;
      authDomain?: string;
      environment?: string;
    };
  }
}

@Injectable({
  providedIn: 'root',
})
export class EnvService {
  private readonly runtimeEnv =
    typeof window !== 'undefined' && window.__env ? window.__env : {};

  get apiUrl(): string {
    return this.runtimeEnv.apiUrl || environment.apiUrl;
  }

  get authDomain(): string {
    return this.runtimeEnv.authDomain || environment.authDomain || '';
  }

  get environment(): string {
    return (
      this.runtimeEnv.environment || environment.environment || 'development'
    );
  }
}
