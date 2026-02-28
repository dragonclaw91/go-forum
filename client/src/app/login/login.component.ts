import { Component, ChangeDetectionStrategy, signal, NgModule, computed, inject, OnInit } from '@angular/core';
import { RouterOutlet, Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormsModule, Validators } from '@angular/forms';
import { AuthService } from '../auth/auth.service';
import { ReactiveFormsModule } from '@angular/forms';

/* we are using enums for type saftey and to avoid passing around magic strings 
and we are defineing it up top to remind ourselves that this is going to be used else where in the app */
export enum AuthAction {
  Login = 'login',
  Signup = 'signup'
}

export interface AuthRequest {
  destination: AuthAction;    
  name: string;
  password: string;     
}


@Component({
  selector: 'app-root',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule
  ],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})



export class LoginComponent implements OnInit {
  signupForm!: FormGroup;
  private authService = inject(AuthService);
  private router = inject(Router)
  private errorTimer: any;



  ngOnInit(): void {
    console.log('Component is now on the DOM!');
  }
  hidePassword = signal(false)
  /* we are using an object here because its easier to keep track of things in the object 
instead of having to remeber to update everything when changes are made
  */
  authState = signal({
    mode: AuthAction.Login,
    label: 'Login'
  });
  isLoginMode = signal(true);
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);
  shakeTrigger = signal(false);
  errorLabel = signal(false);
  username = signal('');
  password = signal('');

  // submitLabel = computed(() =>
  //   this.isLoginMode() ? 'Login' : 'Create Account'
  // );

//    currentState  =  this.authState.update(current => 
//   current.mode === AuthAction.Login 
//     ? { mode: AuthAction.Signup, label: 'Create Account' }
//     : { mode: AuthAction.Login, label: 'Sign In' }
// );

//   determineState(){
//     this.authState.update(current => 
//   current.mode === AuthAction.Login 
//     ? { mode: AuthAction.Signup, label: 'Create Account' }
//     : { mode: AuthAction.Login, label: 'Sign In' }
// );
  // }

  clearError() {
    this.errorMessage.set(null);
    this.errorLabel.set(false)
  }

  togglePassword() {
    this.hidePassword.update(v => !v)
  }

  toggleMode() {
      this.authState.update(current => 
  current.mode === AuthAction.Login 
    ? { mode: AuthAction.Signup, label: 'Create Account' }
    : { mode: AuthAction.Login, label: 'Sign In' }
);
  }

  onInputClear() {
    if (this.errorMessage()) {
      this.clearError();
    }
  }
  // shake the login card for a few seconds
  triggerShake() {
    this.shakeTrigger.set(true);
    this.errorLabel.set(true);
    setTimeout(() => {
      this.shakeTrigger.set(false);
    }, 400);
  }



  onSignIn() {
    // const loginData = { name: this.username(), password: this.password(), destination:this.authState().mode };
   const loginData: AuthRequest = { name: this.username(), password: this.password(), destination:this.authState().mode }
    console.log("login data", loginData.name)
    this.isLoading.set(true);
    this.authService.auth(loginData).subscribe({
      next: (response) => {
        this.isLoading.set(false);
        console.log('Login Successful!', response);
        // this.router.navigate(['/home']);
      },
      error: (err) => {
        console.log("ERR", err)
        const msg = err.error?.error || 'an unexpected error occured'

        this.errorMessage.set(msg);

        /* Kill the old timer
          if we dont problems will happen when users
          try to submit again before the error clears */
        if (this.errorTimer) clearTimeout(this.errorTimer);

        const isMobile = window.innerWidth < 1000;

        if (isMobile) {
          this.errorTimer = setTimeout(() => this.errorMessage.set(null), 5000);
        }

        this.isLoading.set(false);
        this.triggerShake();
      }
    });
  }

}
