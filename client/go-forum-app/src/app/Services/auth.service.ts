import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private apiUrl = 'http://localhost:5000/v1/signup' // Your Go endpoint

  constructor(private http: HttpClient) {}

  login(credentials: any): Observable<any> {
    // This sends a POST request to Go with the user's data
    return this.http.post(this.apiUrl, credentials);
  }
}