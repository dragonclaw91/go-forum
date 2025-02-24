import {ApplicationConfig} from '@angular/core';
import {provideRouter} from '@angular/router';
import {routes} from './app.routes';
import { providePrimeNG } from 'primeng/config';

import customPreset from './mypreset';

export const appConfig: ApplicationConfig = {
  providers: [provideRouter(routes),    providePrimeNG({ theme: { preset: customPreset} })],
  
};
