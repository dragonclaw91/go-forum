import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async'; // Import your routes

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes), provideAnimationsAsync()  // Use the routes here
  ]
})
.catch(err => console.error(err));
