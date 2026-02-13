import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { UserService } from './user.service';

interface LoginResponse {
    access_token: string; // The "Label"
}
@Injectable({ providedIn: 'root' })
export class AuthService {
    private apiUrl = 'http://localhost:5000/v1/auth/login' // Your Go endpoint

    constructor(private http: HttpClient,
                private userService: UserService
                ) { }


    login(credentials: any): Observable<any> {
        // This sends a POST request to Go with the user's data
        return this.http.post<LoginResponse>(this.apiUrl, credentials, { withCredentials: true }).pipe(
            tap(response => {
                console.log("response", response)
                //similar to a hook used to manage side effects
                if (response && response.access_token) {
                    localStorage.setItem('access_token', response.access_token)
// 2. Now call the Sibling Service directly
      this.userService.getUserSettings().subscribe({
        next: (userData) => {
          console.log("Token verified! User is:", userData);
        }
      });
                    this.userService
                    console.log('Access Token saved to LocalStorage')
                }
            })
        );

    }
}