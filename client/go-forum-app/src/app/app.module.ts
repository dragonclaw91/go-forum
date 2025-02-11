import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NbThemeModule, NbLayoutModule, NbCardModule } from '@nebular/theme';
import { AppComponent } from './app.component';

@NgModule({
  declarations: [

    // Other components can be added here
  ],
  imports: [
    BrowserModule,
    NbThemeModule.forRoot(), // Nebular theme setup
    NbLayoutModule,          // Layout module
    NbCardModule,           // Card module
    // Other modules can be added here
  ],
  providers: [],
//   bootstrap: [AppComponent] // Bootstrap your root component
})
export class AppModule { }