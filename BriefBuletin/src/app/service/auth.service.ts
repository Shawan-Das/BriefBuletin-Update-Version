import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  // Create new user account
  createUser(userData: { email: string; password: string; phone: string; userName: string; role?: string }): Observable<any> {
    return this.http.post<any>(this.apiUrl + 'api/auth/create', userData);
  }

  // Create a new admin (wrapper around createUser with role=ADMIN)
  createAdmin(userData: { email: string; password: string; phone: string; userName: string }): Observable<any> {
    const payload = { ...userData, role: 'ADMIN' };
    return this.createUser(payload);
  }

  // Verify OTP
  verifyOtp(login: string, otp: string): Observable<any> {
    return this.http.post<any>(this.apiUrl + 'api/verify-user', { login, otp });
  }

  // Send OTP
  sendOtp(login: string): Observable<any> {
    return this.http.post<any>(this.apiUrl + 'api/send-otp', { login });
  }

  // Login user
  login(username: string, password: string): Observable<any> {
    return this.http.post<any>(this.apiUrl + 'api/auth/login', { login: username, pwd: password }).pipe(
      tap(response => {
        if (response.statusCode === 200 && response.payload?.token) {
          const token = response.payload.token;
          const role = response.payload.role;
          const user_name = response.payload.user_name;
          sessionStorage.setItem('token', token);
          sessionStorage.setItem('role', role);
          sessionStorage.setItem('user_name', user_name);
        }
      })
    );
  }

  // Reset password
  resetPassword(email: string, newPwd: string): Observable<any> {
    return this.http.post<any>(this.apiUrl + 'api/auth/resetpwd', { email, newPwd });
  }

  logout(): void {
    sessionStorage.removeItem('token');
  }

  isLoggedIn(): boolean {
    return !!sessionStorage.getItem('token');
  }
}

