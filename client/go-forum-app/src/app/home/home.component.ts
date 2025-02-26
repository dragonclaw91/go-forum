import { Component } from '@angular/core';
import { MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatSelectModule} from '@angular/material/select';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import { SelectModule } from 'primeng/select';
import { FormsModule } from '@angular/forms'; 
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { CardModule } from 'primeng/card';


@Component({
	selector: 'app-root',
	imports: [  
    CardModule,
    InputIconModule,
    IconFieldModule ,
    FormsModule,
 SelectModule ,
    MatSlideToggleModule,
     MatSelectModule, 
     MatFormFieldModule, 
     MatInputModule 
    ],
  templateUrl: './home.component.html',
  styles:` @use '@angular/material' as mat;
  .my-dropdown-panel {
    @include mat.form-field-overrides((
      outlined-focus-outline-color: orange;
    filled-focus-active-indicator-color: red;
  ));
  }
  `,
  styleUrl: './home.component.css'
})
export class HomeComponent {
  cities: any[];
  selectedCity: any;

  constructor() {
    this.cities = [
        {name: 'New York', code: 'NY'},
        {name: 'London', code: 'LDN'},
        {name: 'Paris', code: 'PRS'}
    ];
}
}
