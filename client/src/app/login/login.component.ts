import { Component, ChangeDetectionStrategy, signal, NgModule, computed, inject, OnInit } from '@angular/core';
import { RouterOutlet, Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormsModule, Validators } from '@angular/forms';
import { AuthService } from '../auth/auth.service';
import { ReactiveFormsModule } from '@angular/forms';





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

  isLoginMode = signal(true);
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);
  shakeTrigger = signal(false);
  errorLabel = signal(false);
  username = signal('');
  password = signal('');

  submitLabel = computed(() =>
    this.isLoginMode() ? 'Login' : 'Create Account'
  );

  clearError() {
    this.errorMessage.set(null);
    this.errorLabel.set(false)
  }

  toggleMode() {
    this.isLoginMode.update(v => !v);
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


//TODO: the signal is not tracking the data 
  onSignIn() {
    const loginData = { name: this.username(), password: this.password()};
    console.log("login data",loginData.name)
    this.isLoading.set(true);
    this.authService.login(loginData).subscribe({
      next: (response) => {
        this.isLoading.set(false);
        console.log('Login Successful!', response);
        // this.router.navigate(['/home']);
      },
      error: (err) => {
        console.log("ERR",err)
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
