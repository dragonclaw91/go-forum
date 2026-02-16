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
  connectionStatus = signal('Initializing...');
  ngOnInit(): void {
    console.log('Component is now on the DOM!');
    // this.checkBackendConnection();
  }

  isLoginMode = signal(true);
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);
  shakeTrigger = signal(false);
  username = signal('')
  password = signal('')
  // Computed signal for the submit button text
  submitLabel = computed(() =>
    this.isLoginMode() ? 'Login' : 'Create Account'
  );



  isVisible = false;



  // toggleVisibility() {
  //   this.isVisible = !this.isVisible
  //   this.onSignIn()

  // }
  clearError(){
      this.errorMessage.set(null);
  }

  toggleMode() {
    this.isLoginMode.update(v => !v);
    this.errorMessage.set(null); // Clear errors when switching modes
  }

  // onSubmit() {

  //     // this.triggerEffect(messageFromJson);

  //     // Replace 'formInvalid' with your actual validation check
  //   const formInvalid = true; 

  //   if (formInvalid) {
  //     this.triggerShake();
  //     return;
  //   }
  //     // Logic for Go backend integration goes here tomorrow
  //   }

  triggerShake() {
    // 1. Turn it on
    this.shakeTrigger.set(true);

    // 2. Turn it off after the animation finishes (400ms)
    // This allows the class to be re-added on the next click
    setTimeout(() => {
      this.shakeTrigger.set(false);
    }, 400);
  }

  onSignIn() {
    const  loginData = { username: '', password: '' };
    this.isLoading.set(true);
    console.log('Sign In:', { username: loginData.username, password: loginData.password });
    this.authService.login(loginData).subscribe({
      next: (response) => {
        this.isLoading.set(false);
        console.log('Login Successful!', response);
        // this.router.navigate(['/home']);
        // Here is where you would redirect to the Reddit feed
      },
      error: (err) => {
        console.log("Response", err.error)
        const msg = err.srror?.error || 'an unexpected error occured'
        this.errorMessage.set(msg);
        setTimeout(() => this.errorMessage.set(null), 5000);
        this.isLoading.set(false);
        this.triggerShake();
        console.error('Login Failed', err);
      }
    });
  }

  // onSignIn() {
  //   console.log('Sign In:', { email: this.loginData.username, password: this.loginData.password });
  //   this.router.navigate(['/home']);
  // }


}
