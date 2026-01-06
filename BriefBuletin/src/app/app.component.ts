import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'BriefBuliteen';
  isAdmin(): boolean {
    const r = sessionStorage.getItem('role') || '';
    return r.toUpperCase() === 'ADMIN';
  }
}
