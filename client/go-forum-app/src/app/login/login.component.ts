import { Component,ChangeDetectionStrategy, } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import {FormsModule} from '@angular/forms';
import { AuthService } from '../Services/auth.service'



@Component({
  selector: 'app-root',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})

export class LoginComponent {

  constructor(private router: Router, private  authService: AuthService) { }
  

loginData = { username: '', password: '' };

isVisible = false;

toggleVisibility() {
  this.isVisible = !this.isVisible
  this.onSignIn()

}



  onSignIn() {
    console.log('Sign In:', { username: this.loginData.username, password: this.loginData.password });
    this.authService.login(this.loginData).subscribe({
      next: (response) => {
        console.log('Login Successful!', response);
        this.router.navigate(['/home']);
        // Here is where you would redirect to the Reddit feed
      },
      error: (err) => {
        console.error('Login Failed', err);
      }
    });
  }

// onSignIn() {
//   console.log('Sign In:', { email: this.loginData.username, password: this.loginData.password });
//   this.router.navigate(['/home']);
// }


}
