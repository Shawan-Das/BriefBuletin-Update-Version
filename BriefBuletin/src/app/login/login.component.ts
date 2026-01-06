import { Component } from '@angular/core';
import { AuthService } from '../service/auth.service';
import { Router } from '@angular/router';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  // View states
  currentView: 'signin' | 'signup' | 'otp' | 'forgot-password' | 'reset-password' = 'signin';
  
  // Sign in fields
  loginEmail: string = '';
  loginPassword: string = '';
  showPassword = false;

  // Sign up fields
  signupEmail: string = '';
  signupPassword: string = '';
  signupConfirmPassword: string = '';
  signupPhone: string = '';
  signupUserName: string = '';
  showSignupPassword = false;

  // OTP fields
  otp: string = '';
  otpEmail: string = '';
  otpPurpose: 'signup' | 'login' | 'forgot-password' = 'signup';
  otpTimer: number = 300; // 5 minutes in seconds
  otpTimerInterval: any;

  // Reset password fields
  resetEmail: string = '';
  resetNewPassword: string = '';
  resetConfirmPassword: string = '';
  showResetPassword = false;

  constructor(private authService: AuthService, private router: Router) {}

  // Toggle between sign in and sign up
  switchView(view: 'signin' | 'signup') {
    this.currentView = view;
    this.resetForms();
  }

  // Reset all forms
  resetForms() {
    this.loginEmail = '';
    this.loginPassword = '';
    this.signupEmail = '';
    this.signupPassword = '';
    this.signupConfirmPassword = '';
    this.signupPhone = '';
    this.signupUserName = '';
    this.otp = '';
    this.otpEmail = '';
    this.resetEmail = '';
    this.resetNewPassword = '';
    this.resetConfirmPassword = '';
    this.showPassword = false;
    this.showSignupPassword = false;
    this.showResetPassword = false;
    this.clearOtpTimer();
  }

  // Sign up
  signup() {
    if (!this.signupEmail || !this.signupPassword || !this.signupPhone || !this.signupUserName) {
      Swal.fire('Validation Error', 'Please fill all required fields', 'error');
      return;
    }

    if (this.signupPassword !== this.signupConfirmPassword) {
      Swal.fire('Validation Error', 'Passwords do not match', 'error');
      return;
    }

    if (this.signupPassword.length < 8) {
      Swal.fire('Validation Error', 'Password must be at least 8 characters long', 'error');
      return;
    }

    const userData = {
      email: this.signupEmail,
      password: this.signupPassword,
      phone: this.signupPhone,
      userName: this.signupUserName
    };

    this.authService.createUser(userData).subscribe(
      (response) => {
        if (response.isSuccess) {
          this.otpEmail = this.signupEmail;
          this.otpPurpose = 'signup';
          this.currentView = 'otp';
          this.startOtpTimer();
          Swal.fire('Success', response.serviceMessage || 'OTP sent to your email', 'success');
        } else {
          Swal.fire('Error', response.serviceMessage || 'Failed to create account', 'error');
        }
      },
      (error) => {
        Swal.fire('Error', error.error?.serviceMessage || 'Failed to create account', 'error');
      }
    );
  }

  loginAsGuest(){
    // this.loginEmail= "BriefBulletin";
    // this.loginPassword= "Guest@000";
    this.authService.login("BriefBulletin", "Guest@000").subscribe(
      (response) => {
        if (response.statusCode === 200) {
          // Login successful, get categories
          this.router.navigate(['/home']);
        }
      }
    );
  }
  // Sign in
  login() {
    if (!this.loginEmail || !this.loginPassword) {
      Swal.fire('Validation Error', 'Please enter email and password', 'error');
      return;
    }

    this.authService.login(this.loginEmail, this.loginPassword).subscribe(
      (response) => {
        if (response.statusCode === 200) {
          // Login successful, get categories
          this.router.navigate(['/home']);
        }
      },
      (error) => {
        // Check if it's a 403 response (email not verified)
        if (error.status === 403 || error.error?.statusCode === 403) {
          // Email not verified, send OTP
          this.otpEmail = this.loginEmail;
          this.otpPurpose = 'login';
          this.currentView = 'otp';
          this.startOtpTimer();
          Swal.fire('Verification Required', error.error?.serviceMessage || 'OTP sent to your email. Please verify yourself to login.', 'info');
        } else {
          // Other errors
          Swal.fire('Login Failed', error.error?.serviceMessage || 'Invalid credentials', 'error');
        }
      }
    );
  }

  // Verify OTP
  verifyOtp() {
    if (!this.otp || this.otp.length !== 6) {
      Swal.fire('Validation Error', 'Please enter a valid 6-digit OTP', 'error');
      return;
    }

    // Ensure OTP contains only numbers
    if (!/^\d+$/.test(this.otp)) {
      Swal.fire('Validation Error', 'OTP must contain only numbers', 'error');
      return;
    }

    this.authService.verifyOtp(this.otpEmail, this.otp).subscribe(
      (response) => {
        if (response.statusCode === 200) {
          this.clearOtpTimer();
          if (this.otpPurpose === 'signup') {
            Swal.fire('Success', 'Account created successfully! Please login.', 'success').then(() => {
              this.currentView = 'signin';
              this.resetForms();
            });
          } else if (this.otpPurpose === 'login') {
            // After OTP verification for login, try login again
            this.authService.login(this.otpEmail, this.loginPassword).subscribe(
              (loginResponse) => {
                if (loginResponse.statusCode === 200) {
                  this.router.navigate(['/home']);
                }
              },
              () => {
                Swal.fire('Error', 'Login failed after verification', 'error');
              }
            );
          } else if (this.otpPurpose === 'forgot-password') {
            // Move to reset password view
            this.currentView = 'reset-password';
            this.otp = '';
          }
        } else {
          Swal.fire('Invalid OTP', 'The OTP you entered is invalid or expired', 'error');
        }
      },
      (error) => {
        Swal.fire('Error', error.error?.serviceMessage || 'Invalid OTP', 'error');
      }
    );
  }

  // Resend OTP
  resendOtp() {
    // Prevent resending if timer is still active
    if (this.otpTimer > 0) {
      return;
    }

    this.authService.sendOtp(this.otpEmail).subscribe(
      (response) => {
        this.startOtpTimer();
        Swal.fire('Success', 'OTP has been resent to your email', 'success');
      },
      (error) => {
        Swal.fire('Error', error.error?.serviceMessage || 'Failed to resend OTP', 'error');
      }
    );
  }

  // Forgot password - send OTP
  forgotPassword() {
    if (!this.resetEmail) {
      Swal.fire('Validation Error', 'Please enter your email address', 'error');
      return;
    }

    this.authService.sendOtp(this.resetEmail).subscribe(
      (response) => {
        this.otpEmail = this.resetEmail;
        this.otpPurpose = 'forgot-password';
        this.currentView = 'otp';
        this.startOtpTimer();
        Swal.fire('Success', 'OTP sent to your email', 'success');
      },
      (error) => {
        Swal.fire('Error', error.error?.serviceMessage || 'Failed to send OTP', 'error');
      }
    );
  }

  // Reset password
  resetPassword() {
    if (!this.resetNewPassword || !this.resetConfirmPassword) {
      Swal.fire('Validation Error', 'Please enter and confirm your new password', 'error');
      return;
    }

    if (this.resetNewPassword !== this.resetConfirmPassword) {
      Swal.fire('Validation Error', 'Passwords do not match', 'error');
      return;
    }

    if (this.resetNewPassword.length < 8) {
      Swal.fire('Validation Error', 'Password must be at least 8 characters long', 'error');
      return;
    }

    this.authService.resetPassword(this.resetEmail, this.resetNewPassword).subscribe(
      (response) => {
        if (response.statusCode === 200 || response.isSuccess) {
          Swal.fire('Success', 'Password reset successfully! Please login.', 'success').then(() => {
            this.currentView = 'signin';
            this.resetForms();
          });
        } else {
          Swal.fire('Error', response.serviceMessage || 'Failed to reset password', 'error');
        }
      },
      (error) => {
        Swal.fire('Error', error.error?.serviceMessage || 'Failed to reset password', 'error');
      }
    );
  }

  // OTP Timer
  startOtpTimer() {
    this.otpTimer = 300; // 5 minutes
    this.clearOtpTimer();
    this.otpTimerInterval = setInterval(() => {
      this.otpTimer--;
      if (this.otpTimer <= 0) {
        this.clearOtpTimer();
      }
    }, 1000);
  }

  clearOtpTimer() {
    if (this.otpTimerInterval) {
      clearInterval(this.otpTimerInterval);
      this.otpTimerInterval = null;
    }
  }

  getOtpTimerDisplay(): string {
    const minutes = Math.floor(this.otpTimer / 60);
    const seconds = this.otpTimer % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  }

  // Password visibility toggles
  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  toggleSignupPasswordVisibility() {
    this.showSignupPassword = !this.showSignupPassword;
  }

  toggleResetPasswordVisibility() {
    this.showResetPassword = !this.showResetPassword;
  }

  // Navigate back
  goBack() {
    if (this.currentView === 'otp') {
      if (this.otpPurpose === 'forgot-password') {
        this.currentView = 'forgot-password';
      } else {
        this.currentView = this.otpPurpose === 'signup' ? 'signup' : 'signin';
      }
      this.clearOtpTimer();
    } else if (this.currentView === 'reset-password') {
      this.currentView = 'otp';
    } else if (this.currentView === 'forgot-password') {
      this.currentView = 'signin';
    }
    this.otp = '';
  }
}
