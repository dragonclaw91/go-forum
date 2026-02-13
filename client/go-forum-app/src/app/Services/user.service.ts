import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(private http: HttpClient) { }

  // This is the "Protected" request we discussed
  getUserSettings() {
    const token = localStorage.getItem('access_token');
    
    return this.http.get('http://localhost:8080/user/settings', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
  }
}