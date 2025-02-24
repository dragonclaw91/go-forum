import { Component,ChangeDetectionStrategy, } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import {FormsModule} from '@angular/forms';



@Component({
  selector: 'app-root',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})

export class LoginComponent {

  constructor(private router: Router) { }

  email = '';
  password = '';

isVisible = false;

toggleVisibility() {
  this.isVisible = !this.isVisible

}
onSignIn() {
  console.log('Sign In:', { email: this.email, password: this.password });
  this.router.navigate(['/home']);
}


}
