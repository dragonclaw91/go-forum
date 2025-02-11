import { Component,ChangeDetectionStrategy, } from '@angular/core';
import { RouterOutlet } from '@angular/router';
// import { LoginComponent } from './login/login.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
template: `<router-outlet>`,

})
export class AppComponent {

}
