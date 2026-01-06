import { Component, OnInit } from '@angular/core';
import { AuthService } from '../service/auth.service';
import { Router } from '@angular/router';
import { Item } from '../models/item.model';
import { ApiServiceService } from '../service/api-service.service';
import Swal from 'sweetalert2';

interface Article {
  id: number;
  title: string;
  summary: string;
  content: string;
  featured_image: string;
  category_id: number;
  status: string;
  published_at: string;
  views_count: number;
  created_at: string;
  updated_at: string;
  source_url: string;
}

interface Comment {
  id?: number;
  article_id: number;
  user_name: string;
  user_email: string;
  content: string;
  created_at?: string;
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {
  userRole : string = '';
  userName : string = '';
  isAdmin: boolean = false;
  articles: Article[] = [];
  currentPage: number = 1;
  isLoading: boolean = false;
  hasMoreArticles: boolean = true;
  isDraftingArticle: boolean = false;
  selectedLanguage: string = 'en'; // Default language is English
  totalNews: number = 0;

  constructor(
    private authService: AuthService,
    private router: Router,
    private apiService: ApiServiceService
  ) {}

  ngOnInit() {
    this.userRole = sessionStorage.getItem('role') || '';
    this.userName = sessionStorage.getItem('user_name') || '';
    this.isAdmin = this.userRole === 'ADMIN';
    this.loadArticles(1);
  }

  draftArticle(articleId: number) {
    if (this.isDraftingArticle) return;
    
    // Show confirmation dialog
    Swal.fire({
      title: 'Move to Draft?',
      text: 'This article will be moved to draft and removed from the public view.',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc3545',
      cancelButtonColor: '#6c757d',
      confirmButtonText: 'Yes, move to draft',
      cancelButtonText: 'Cancel',
      reverseButtons: true,
      customClass: {
        confirmButton: 'btn-draft-confirm',
        cancelButton: 'btn-draft-cancel'
      }
    }).then((result) => {
      if (result.isConfirmed) {
        this.performDraftArticle(articleId);
      }
    });
  }

  // New method to handle the actual drafting
  private performDraftArticle(articleId: number) {
    this.isDraftingArticle = true;
    
    this.apiService.draftArticle(articleId).subscribe(
      (response) => {
        this.isDraftingArticle = false;
        
        if (response.isSuccess) {
          // Remove the article from the local array without reloading
          this.articles = this.articles.filter(article => article.id !== articleId);
          
          // Show success notification
          Swal.fire({
            title: 'Moved to Draft!',
            text: 'The article has been successfully moved to draft.',
            icon: 'success',
            timer: 1500,
            showConfirmButton: false,
            customClass: {
              popup: 'draft-success-notification'
            }
          });
        } else {
          Swal.fire(
            'Error', 
            response.serviceMessage || 'Failed to move article to draft', 
            'error'
          );
        }
      },
      (error) => {
        this.isDraftingArticle = false;
        Swal.fire(
          'Error', 
          'An error occurred while moving the article to draft', 
          'error'
        );
      }
    );
  }

  loadArticles(page: number) {
    if (this.isLoading) return;

    this.isLoading = true;
    this.totalNews = page === 1 ? 0 : this.articles.length;
    this.apiService.getArticles(page, this.selectedLanguage, this.totalNews).subscribe(
      (response) => {
        if (response.isSuccess && response.payload) {
          if (page === 1) {
            this.articles = response.payload;
          } else {
            // Append new articles to existing ones
            this.articles = [...this.articles, ...response.payload];
          }
          
          // Check if there are more articles (if payload is empty or has fewer items, no more pages)
          if (!response.payload || response.payload.length === 0) {
            this.hasMoreArticles = false;
          }
        } else {
          this.hasMoreArticles = false;
        }
        this.isLoading = false;
      },
      (error) => {
        this.isLoading = false;
        this.hasMoreArticles = false;
        Swal.fire('Error', 'Failed to load articles', 'error');
      }
    );
  }

  switchLanguage(language: string) {
    if (this.selectedLanguage === language) return; // No change needed
    
    this.selectedLanguage = language;
    this.currentPage = 1; // Reset page to 1 when switching language
    this.articles = []; // Clear existing articles
    this.hasMoreArticles = true; // Reset hasMoreArticles flag
    this.loadArticles(1); // Load first page of new language
  }

  loadMore() {
    if (!this.isLoading && this.hasMoreArticles) {
      this.currentPage++;
      this.loadArticles(this.currentPage);
    }
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  formatDateTime(dateString: string): string {
    const date = new Date(dateString);
    const dateStr = date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
    const timeStr = date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true
    });
    return `${dateStr} at ${timeStr}`;
  }

  readNews(articleId: number) {
    // Find the article
    const article = this.articles.find(a => a.id === articleId);
    if (!article) return;

    // Call API to increment view count (silently)
    this.apiService.readNews(articleId).subscribe(
      () => {
        // Update the view count locally
        article.views_count++;
      },
      () => {
        // Silently handle errors - view count increment failed
      }
    );

    // Fetch comments for this article
    this.apiService.getComments(articleId).subscribe(
      (commentsResponse) => {
        const comments = commentsResponse.payload || commentsResponse || [];
        this.showArticleModal(article, comments);
      },
      () => {
        // If comments fail to load, show modal without comments
        this.showArticleModal(article, []);
      }
    );
  }

  showArticleModal(article: Article, comments: any[]) {
    // Generate comments HTML
    let commentsHtml = '';
    if (comments && comments.length > 0) {
      commentsHtml = comments.map((comment: any) => {
        const commentDate = comment.created_at ? this.formatDateTime(comment.created_at) : 'Recently';
        return `
          <div class="comment-item" style="background: #f8f9fa; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; border-left: 3px solid #007bff;">
            <div class="comment-header" style="display: flex; justify-content: space-between; align-items: start; margin-bottom: 0.5rem; flex-wrap: wrap;">
              <div class="comment-author">
                <strong style="color: #212529; font-size: 0.95rem; display: block;">${this.escapeHtml(comment.user_name || 'Anonymous')}</strong>
                <span class="comment-email" style="color: #6c757d; font-size: 0.8rem; display: block; margin-top: 0.25rem;">${comment.user_email || ''}</span>
              </div>
              <span class="comment-date" style="color: #6c757d; font-size: 0.75rem; margin-top: 0.25rem;">${commentDate}</span>
            </div>
            <p class="comment-content" style="color: #333; margin: 0; line-height: 1.6; font-size: 0.9rem; word-wrap: break-word;">${this.escapeHtml(comment.content || '')}</p>
          </div>
        `;
      }).join('');
    } else {
      commentsHtml = '<p style="color: #6c757d; text-align: center; padding: 2rem; font-size: 0.9rem;">No comments yet. Be the first to comment!</p>';
    }

    // Comment form HTML
    const commentFormHtml = `
      <div class="comment-form-wrapper" style="margin-top: 2rem; padding-top: 1.5rem; border-top: 2px solid #dee2e6;">
        <h5 class="comment-form-title" style="margin-bottom: 1rem; color: #212529; font-size: 1rem;"><i class="fas fa-comment me-2"></i>Add a Comment</h5>
        <form id="commentForm" style="margin-bottom: 0;">
          <div style="margin-bottom: 1rem;">
            <input type="text" id="commentUserName" placeholder="Your Name" required 
                   class="comment-input" style="width: 100%; padding: 0.5rem; border: 1px solid #ced4da; border-radius: 4px; font-size: 0.9rem; box-sizing: border-box;" />
          </div>
          <div style="margin-bottom: 1rem;">
            <input type="email" id="commentUserEmail" placeholder="Your Email" required 
                   class="comment-input" style="width: 100%; padding: 0.5rem; border: 1px solid #ced4da; border-radius: 4px; font-size: 0.9rem; box-sizing: border-box;" />
          </div>
          <div style="margin-bottom: 1rem;">
            <textarea id="commentContent" placeholder="Write your comment..." required rows="4"
                      class="comment-textarea" style="width: 100%; padding: 0.5rem; border: 1px solid #ced4da; border-radius: 4px; font-size: 0.9rem; resize: vertical; box-sizing: border-box;"></textarea>
          </div>
          <button type="submit" id="submitCommentBtn" 
                  class="comment-submit-btn" style="background: #007bff; color: white; border: none; padding: 0.6rem 1.5rem; border-radius: 4px; cursor: pointer; font-weight: 500; width: 100%; font-size: 0.9rem;">
            <i class="fas fa-paper-plane me-1"></i>Post Comment
          </button>
        </form>
      </div>
    `;

    Swal.fire({
      title: article.title,
      html: `
        <div style="text-align: left;" class="article-modal-wrapper">
          <div class="article-content-section" style="margin-bottom: 2rem;">
            ${article.featured_image ? `<img src="${article.featured_image}" class="article-modal-image" style="width: 100%; max-height: 500px; object-fit: contain; border-radius: 8px; margin-bottom: 1rem; display: block;" onerror="this.onerror=null; this.src='${this.getPlaceholderImage()}'" />` : `<img src="${this.getPlaceholderImage()}" class="article-modal-image" style="width: 100%; max-height: 500px; object-fit: contain; border-radius: 8px; margin-bottom: 1rem; display: block;" />`}
            <p class="article-meta-info" style="color: #6c757d; margin-bottom: 1rem; font-size: 0.9rem;">
              <i class="far fa-calendar-alt me-1"></i><strong>Published:</strong> ${this.formatDateTime(article.published_at)}
              <span class="meta-separator ms-2 ms-md-3"><i class="far fa-eye me-1"></i><strong>Views:</strong> ${article.views_count}</span>
            </p>
            <div class="article-body" style="line-height: 1.8; color: #333; white-space: pre-wrap; margin-bottom: 1.5rem; font-size: 0.95rem;">${article.content || article.summary || 'No content available.'}</div>
            ${article.source_url ? `<div class="source-link" style="margin-bottom: 1.5rem; padding-bottom: 1rem; border-bottom: 1px solid #dee2e6;"><a href="${article.source_url}" target="_blank" style="color: #007bff; text-decoration: none; font-size: 0.9rem;"><i class="fas fa-external-link-alt me-1"></i>View Original Source</a></div>` : ''}
          </div>
          
          <div class="comments-section" style="border-top: 2px solid #dee2e6; padding-top: 1.5rem; margin-top: 2rem;">
            <h5 class="comments-header" style="margin-bottom: 1rem; color: #212529; font-size: 1.1rem;"><i class="fas fa-comments me-2"></i>Comments (${comments.length})</h5>
            <div id="commentsContainer" class="comments-list" style="max-height: 40vh; overflow-y: auto; margin-bottom: 1.5rem;">
              ${commentsHtml}
            </div>
            ${commentFormHtml}
          </div>
        </div>
      `,
      width: '90%',
      showCloseButton: true,
      showConfirmButton: false,
      didOpen: () => {
        // Scroll to top to focus on article content
        const modalContent = document.querySelector('.swal2-html-container');
        if (modalContent) {
          modalContent.scrollTop = 0;
        }
        
        // Handle comment form submission
        const form = document.getElementById('commentForm') as HTMLFormElement;
        const submitBtn = document.getElementById('submitCommentBtn') as HTMLButtonElement;
        
        if (form) {
          form.addEventListener('submit', (e) => {
            e.preventDefault();
            const userName = (document.getElementById('commentUserName') as HTMLInputElement)?.value;
            const userEmail = (document.getElementById('commentUserEmail') as HTMLInputElement)?.value;
            const content = (document.getElementById('commentContent') as HTMLTextAreaElement)?.value;

            if (!userName || !userEmail || !content) {
              Swal.fire('Validation Error', 'Please fill all fields', 'error');
              return;
            }

            // Disable submit button
            if (submitBtn) {
              submitBtn.disabled = true;
              submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin me-1"></i>Posting...';
            }

            // Submit comment
            this.apiService.createComment({
              article_id: article.id,
              user_name: userName,
              user_email: userEmail,
              content: content
            }).subscribe(
              (response) => {
                if (response.isSuccess || response.statusCode === 200) {
                  // Reload comments
                  this.apiService.getComments(article.id).subscribe(
                    (commentsResponse) => {
                      const newComments = commentsResponse.payload || commentsResponse || [];
                      this.updateCommentsInModal(newComments);
                      // Reset form
                      form.reset();
                      if (submitBtn) {
                        submitBtn.disabled = false;
                        submitBtn.innerHTML = '<i class="fas fa-paper-plane me-1"></i>Post Comment';
                      }
                      // Swal.fire('Success', 'Comment posted successfully!', 'success');
                    },
                    () => {
                      if (submitBtn) {
                        submitBtn.disabled = false;
                        submitBtn.innerHTML = '<i class="fas fa-paper-plane me-1"></i>Post Comment';
                      }
                    }
                  );
                } else {
                  if (submitBtn) {
                    submitBtn.disabled = false;
                    submitBtn.innerHTML = '<i class="fas fa-paper-plane me-1"></i>Post Comment';
                  }
                  Swal.fire('Error', response.serviceMessage || 'Failed to post comment', 'error');
                }
              },
              (error) => {
                if (submitBtn) {
                  submitBtn.disabled = false;
                  submitBtn.innerHTML = '<i class="fas fa-paper-plane me-1"></i>Post Comment';
                }
                Swal.fire('Error', error.error?.serviceMessage || 'Failed to post comment', 'error');
              }
            );
          });
        }
      },
      customClass: {
        popup: 'article-modal',
        htmlContainer: 'article-modal-content'
      }
    });
  }

  updateCommentsInModal(comments: any[]) {
    const commentsContainer = document.getElementById('commentsContainer');
    if (!commentsContainer) return;

    let commentsHtml = '';
    if (comments && comments.length > 0) {
      commentsHtml = comments.map((comment: any) => {
        const commentDate = comment.created_at ? this.formatDateTime(comment.created_at) : 'Recently';
        return `
          <div class="comment-item" style="background: #f8f9fa; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; border-left: 3px solid #007bff;">
            <div class="comment-header" style="display: flex; justify-content: space-between; align-items: start; margin-bottom: 0.5rem; flex-wrap: wrap;">
              <div class="comment-author">
                <strong style="color: #212529; font-size: 0.95rem; display: block;">${this.escapeHtml(comment.user_name || 'Anonymous')}</strong>
                <span class="comment-email" style="color: #6c757d; font-size: 0.8rem; display: block; margin-top: 0.25rem;">${comment.user_email || ''}</span>
              </div>
              <span class="comment-date" style="color: #6c757d; font-size: 0.75rem; margin-top: 0.25rem;">${commentDate}</span>
            </div>
            <p class="comment-content" style="color: #333; margin: 0; line-height: 1.6; font-size: 0.9rem; word-wrap: break-word;">${this.escapeHtml(comment.content || '')}</p>
          </div>
        `;
      }).join('');
    } else {
      commentsHtml = '<p style="color: #6c757d; text-align: center; padding: 2rem; font-size: 0.9rem;">No comments yet. Be the first to comment!</p>';
    }

    commentsContainer.innerHTML = commentsHtml;
    
    // Update comment count in header
    const commentHeader = document.querySelector('.comments-header');
    if (commentHeader) {
      commentHeader.innerHTML = `<i class="fas fa-comments me-2"></i>Comments (${comments.length})`;
    }
  }

  escapeHtml(text: string): string {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  truncateText(text: string, maxLength: number = 150): string {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
  }

  onImageError(event: Event) {
    const img = event.target as HTMLImageElement;
    if (img) {
      // Use SVG data URI as placeholder instead of external URL
      img.src = this.getPlaceholderImage();
    }
  }

  getPlaceholderImage(): string {
    // Return SVG data URI as placeholder (simple gray box with text)
    const svg = `<svg width="800" height="400" xmlns="http://www.w3.org/2000/svg"><rect width="800" height="400" fill="#f0f0f0"/><text x="50%" y="50%" font-family="Arial" font-size="24" fill="#999" text-anchor="middle" dy=".3em">No Image</text></svg>`;
    return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg);
  }

  logout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}