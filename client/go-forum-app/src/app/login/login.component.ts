import { Component, ChangeDetectionStrategy, signal,  NgModule, computed  } from '@angular/core';
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

export class LoginComponent {
  signupForm!: FormGroup;
  // This simulates the JSON object coming from Go
mockBackendError = {
  error: {
    error: "That username has already been claimed by another traveler."
  }
};
  
  constructor(private router: Router, private authService: AuthService, private fb: FormBuilder) {
    this.signupForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required]
    });

  }
  
  isLoginMode = signal(true);
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);
  shakeTrigger = signal(false);
// Computed signal for the submit button text
  submitLabel = computed(() => 
    this.isLoginMode() ? 'Login' : 'Create Account'
  );

  loginData = { username: '', password: '' };

  isVisible = false;

  

  // toggleVisibility() {
  //   this.isVisible = !this.isVisible
  //   this.onSignIn()

  // }

toggleMode() {
    this.isLoginMode.update(v => !v);
    this.errorMessage.set(null); // Clear errors when switching modes
  }

onSubmit() {
    this.isLoading.set(true);
    // this.triggerEffect(messageFromJson);

    // Replace 'formInvalid' with your actual validation check
  const formInvalid = true; 

  if (formInvalid) {
    this.triggerShake();
    return;
  }
    // Logic for Go backend integration goes here tomorrow
  }

  triggerShake() {
  // 1. Turn it on
  this.shakeTrigger.set(true);

  // 2. Turn it off after the animation finishes (400ms)
  // This allows the class to be re-added on the next click
  setTimeout(() => {
    this.shakeTrigger.set(false);
  }, 400);
}

  // onSignIn() {
  //   console.log('Sign In:', { username: this.loginData.username, password: this.loginData.password });
  //   this.authService.login(this.loginData).subscribe({
  //     next: (response) => {
  //       console.log('Login Successful!', response);
  //       this.router.navigate(['/home']);
  //       // Here is where you would redirect to the Reddit feed
  //     },
  //     error: (err) => {
  //       console.error('Login Failed', err);
  //     }
  //   });
  // }

  // onSignIn() {
  //   console.log('Sign In:', { email: this.loginData.username, password: this.loginData.password });
  //   this.router.navigate(['/home']);
  // }


}
