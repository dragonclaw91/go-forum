import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { BsDropdownModule } from 'ngx-bootstrap/dropdown';



@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, BsDropdownModule],
  template: `<router-outlet></router-outlet>`,

})
export class AppComponent {

}
