import { ApplicationConfig } from '@angular/core';
// import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { providePrimeNG } from 'primeng/config';
import Lara from '@primeng/themes/lara'; // This is the file you found!

export const appConfig: ApplicationConfig = {
  providers: [
    // provideAnimationsAsync(),
    providePrimeNG({
        theme: {
            preset: Lara,
            options: {
                darkModeSelector: '.my-app-dark' // Optional: stops it from forcing dark mode
            }
        }
    })
  ]
};