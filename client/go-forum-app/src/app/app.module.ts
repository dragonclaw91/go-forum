import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BsDropdownModule } from 'ngx-bootstrap/dropdown';
import { AppComponent } from './app.component';

@NgModule({
  imports: [
    BrowserModule,
    BsDropdownModule.forRoot() // Ensure this line is present
  ],
  providers: [],
  bootstrap: []  // This should be correctly referencing AppComponent
})
export class AppModule { }