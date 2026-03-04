import { Component, ChangeDetectionStrategy, signal, NgModule, computed, inject, OnInit } from '@angular/core';
import { RouterOutlet, Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormsModule, Validators } from '@angular/forms';
import { AuthService } from '../auth/auth.service';
import { ReactiveFormsModule } from '@angular/forms';
import { addIcons } from 'ionicons';
import { eyeOutline, eyeOffOutline } from 'ionicons/icons';
import { IonButton, IonIcon, IonInput, IonItem, IonList } from '@ionic/angular/standalone';

/* we are using enums for type saftey and to avoid passing around magic strings 
and we are defineing it up top to remind ourselves that this is going to be used else where in the app */
export enum AuthAction {
  Login = 'login',
  Signup = 'signup'
}
export interface AuthState {
  mode: AuthAction;
  isLoading: boolean;
  errorMessage: string | null;
  shakeTrigger: boolean;
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
    ReactiveFormsModule,
    IonButton, 
    IonIcon, 
    IonInput,
     IonItem, 
     IonList
  ],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})



export class LoginComponent implements OnInit {
  signupForm!: FormGroup;
  private authService = inject(AuthService);
  private router = inject(Router)
  private errorTimer: any;




// simialar to use effect in react
  ngOnInit(): void {
    console.log('Component is now on the DOM!');
    addIcons({ 'eye-outline': eyeOutline, 'eye-off-outline': eyeOffOutline });
  }
  hidePassword = signal(false)
  /* we are using an object here because its easier to keep track of things in the object 
instead of having to remeber to update everything when changes are made
  */
  authState = signal<AuthState>({
    mode: AuthAction.Login,
    isLoading: false,
    errorMessage: null,
    shakeTrigger: false,
    name: '',
    password: '',
    destination: AuthAction.Login
  });

  // because the new oblect is defined later we overwrite the existing object
  // basiclly copying the state as is then comparing ot using the spread opeartor
  private patchState(patch: Partial<AuthState>) {
    this.authState.update(current => ({
      ...current,
      ...patch
    }));
  }

  // username = signal('');
  // password = signal('');

  submitLabel = computed(() =>
    this.authState().mode === AuthAction.Login ? 'Login' : 'Create Account'
  );

  footerButton = computed(() => this.authState().mode === AuthAction.Login ? 'Don\'t have an account? Register' : 'Already have an account? Login')

  // using the double not nto prevent things like null to be eevaluated as false
  errorLabel = computed(() => !!this.authState().errorMessage)


  clearError() {
    this.patchState({  errorMessage: null })
  }

  togglePassword() {
    this.hidePassword.update(v => !v)
  }


    togglePasswordType = computed(() => 
    this.hidePassword() ? 'password' : 'text'
  )



  isLoginMode = computed(() => {
    return this.authState().mode === AuthAction.Login
  }
  );



  toggleMode() {
    this.authState().mode === AuthAction.Login ? this.patchState({ mode: AuthAction.Signup }) : this.patchState({ mode: AuthAction.Login })
  }


  // shake the login card for a few seconds
  triggerShake() {
    this.patchState({
      shakeTrigger: true
    });
    setTimeout(() => {
      this.patchState({ shakeTrigger: false })
    }, 400);
  }

 updateField(key: string, value: string) {
  this.authState.update(state => ({
    ...state,
    [key]: value,        
    errorMessage: null   
  }));
}


  onSignIn() {
    const loginData = { name: this.authState().name, password: this.authState().password, destination: this.authState().mode }
    console.log("login data", loginData.name)

    this.patchState({ isLoading: true })
    this.authService.auth(loginData).subscribe({
      next: (response) => {
        // this.isLoading.set(false);
        this.patchState({ isLoading: false })
        console.log('Login Successful!', response);
        // this.router.navigate(['/home']);
      },
      error: (err) => {
        console.log("ERR", err)
        const msg = err.error?.error || 'an unexpected error occured'
        this.patchState({ errorMessage: msg })
        /* Kill the old timer
          if we dont problems will happen when users
          try to submit again before the error clears */
        if (this.errorTimer) clearTimeout(this.errorTimer);
        const isMobile = window.innerWidth < 1000;
        if (isMobile) {
          this.errorTimer = setTimeout(() => this.patchState({ errorMessage: null,  }), 5000);
        }
        this.patchState({ isLoading: false })
        this.triggerShake();
      }
    });
  }

}
