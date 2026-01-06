import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { HomeComponent } from './home/home.component';
import { AuthGuard } from './auth.guard';
import { AdminPanelComponent } from './admin/admin-panel/admin-panel.component';

const routes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'home', component: HomeComponent, canActivate: [AuthGuard] },
  { path: 'admin', component: AdminPanelComponent, canActivate: [AuthGuard] },
  { path: 'admin/create-article', component: AdminPanelComponent, canActivate: [AuthGuard] },
  { path: 'admin/edit-article', component: AdminPanelComponent, canActivate: [AuthGuard] },
  { path: 'admin/approve-article', component: AdminPanelComponent, canActivate: [AuthGuard] },
  { path: 'admin/active-comment', component: AdminPanelComponent, canActivate: [AuthGuard] },
  { path: 'admin/create-admin', component: AdminPanelComponent, canActivate: [AuthGuard] },
  { path: '', redirectTo: '/login', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
