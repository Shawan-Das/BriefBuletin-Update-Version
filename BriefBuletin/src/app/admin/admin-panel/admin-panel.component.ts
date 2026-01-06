import { Component, OnInit } from '@angular/core';
import { ApiServiceService } from '../../service/api-service.service';
import { AuthService } from '../../service/auth.service';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-admin-panel',
  templateUrl: './admin-panel.component.html',
  styleUrls: ['./admin-panel.component.scss']
})
export class AdminPanelComponent implements OnInit {
  role = sessionStorage.getItem('role');
  activeTab: string = 'create-article'; // create-article | edit-article | approve-article | active-comment | create-admin

  // Common
  isLoading = false;
  error = '';
  success = '';

  // Approve Article Data
  drafts: any[] = [];

  // Active Comments Data
  comments: any[] = [];

  // Create Admin Data
  adminForm = {
    name: '',
    email: '',
    password: '',
    phone: ''
  };

  constructor(
    private api: ApiServiceService,
    private auth: AuthService
  ) {}

  ngOnInit(): void {
    if (!this.isAdmin()) {
      this.error = 'You are not authorized to view this page.';
      return;
    }
  }

  isAdmin(): boolean {
    return (this.role || '').toUpperCase() === 'ADMIN';
  }

  switchTab(tab: string) {
    this.activeTab = tab;
    this.error = '';
    this.success = '';

    // Load data for specific tabs
    if (tab === 'approve-article') {
      this.loadDrafts();
    } else if (tab === 'active-comment') {
      this.loadComments();
    }
  }

  // ============= APPROVE ARTICLE SECTION =============
  loadDrafts() {
    this.isLoading = true;
    this.api.getDraftArticles(1).subscribe({
      next: (res) => {
        this.drafts = res.payload || [];
        this.isLoading = false;
      },
      error: () => {
        this.error = 'Failed to load drafts';
        this.isLoading = false;
      }
    });
  }

  approveDraft(articleId: number) {
    this.api.approveArticle(articleId).subscribe({
      next: () => {
        this.drafts = this.drafts.filter(d => d.id !== articleId);
        this.success = 'Article approved successfully!';
        setTimeout(() => this.success = '', 2000);
      },
      error: () => {
        this.error = 'Failed to approve article';
      }
    });
  }

  // ============= ACTIVE COMMENTS SECTION =============
  loadComments() {
    this.isLoading = true;
    this.api.getApprovalDueComments().subscribe({
      next: (res) => {
        this.comments = res.payload || [];
        this.isLoading = false;
      },
      error: () => {
        this.error = 'Failed to load comments';
        this.isLoading = false;
      }
    });
  }

  activateComment(commentId: number) {
    this.api.activateComment(commentId).subscribe({
      next: () => {
        this.comments = this.comments.filter(c => c.comment_id !== commentId);
        this.success = 'Comment activated!';
        setTimeout(() => this.success = '', 2000);
      },
      error: () => {
        this.error = 'Failed to activate comment';
      }
    });
  }

  confirmArchiveComment(commentId: number) {
    Swal.fire({
      title: 'Archive comment?',
      text: 'This will archive (disable) the comment and remove it from the pending list.',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Yes, archive it',
      cancelButtonText: 'Cancel',
      confirmButtonColor: '#dc3545'
    }).then(result => {
      if (result.isConfirmed) {
        this.archiveComment(commentId);
      }
    });
  }

  archiveComment(commentId: number) {
    this.api.disableComment(commentId).subscribe({
      next: () => {
        this.comments = this.comments.filter(c => c.comment_id !== commentId);
        this.success = 'Comment archived!';
        setTimeout(() => this.success = '', 2000);
      },
      error: () => {
        this.error = 'Failed to archive comment';
      }
    });
  }

  // ============= CREATE ADMIN SECTION =============
  submitCreateAdmin() {
    this.error = '';
    this.success = '';

    if (!this.adminForm.name || !this.adminForm.email || !this.adminForm.password) {
      this.error = 'Please fill all required fields (Name, Email, Password).';
      return;
    }

    this.isLoading = true;
    const userData = {
      userName: this.adminForm.name,
      email: this.adminForm.email,
      password: this.adminForm.password,
      phone: this.adminForm.phone
    };

    this.auth.createAdmin(userData).subscribe({
      next: (res) => {
        this.success = 'Admin user created successfully!';
        this.adminForm = { name: '', email: '', password: '', phone: '' };
        this.isLoading = false;
        setTimeout(() => this.success = '', 2000);
      },
      error: (err) => {
        this.error = err?.error?.serviceMessage || 'Failed to create admin user';
        this.isLoading = false;
      }
    });
  }

  resetAdminForm() {
    this.adminForm = { name: '', email: '', password: '', phone: '' };
    this.error = '';
    this.success = '';
  }

  // ============= UTILITY =============
  formatDateTime(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
}
