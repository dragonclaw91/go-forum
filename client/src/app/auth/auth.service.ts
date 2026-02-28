import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { UserService } from '../Services/user.service';

interface LoginResponse {
    access_token: string; // The "Label"
}
@Injectable({ providedIn: 'root' })
export class AuthService {

    private baseUrl = 'http://localhost:5000/v1/auth'

    constructor(private http: HttpClient,
        private userService: UserService
    ) { }




    auth(credentials: any, ): Observable<any> {
        // This sends a POST request to Go with the user's data
        /* we are using intrpolation here because it can get hard to read if we use concat to pull everything together 
        plus it leaves the option to add simple logic if need
        */
        const url = `${this.baseUrl}/${credentials.destination}`
const {destination,...payload} = credentials
        return this.http.post<LoginResponse>(url, credentials, { withCredentials: true }).pipe(

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