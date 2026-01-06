import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { Item } from '../models/item.model';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ApiServiceService {
  private url = environment.apiUrl;

  constructor(private http: HttpClient) {}

  private getToken(): string | null {
    return sessionStorage.getItem('token');
  }

  private getHeaders(): HttpHeaders {
    const token = this.getToken();
    return new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    });
  }

    // --------------- category section ----------------
    // Get categories (after successful login)
    getCategories(): Observable<any> {
      return this.http.get<any>(this.url + 'api/category', {
      });
    }

    // --------------- article section ----------------
    // Get published articles with pagination
    getArticles(page: number = 1, lang: string = 'en', totalNews: number = 0): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/articles?page=${page}&lang=${lang}&totalNews=${totalNews}`, { headers });
    }
    // get draft articles for admin
    getDraftArticles(page: number = 1): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/draft-article-list?page=${page}`, { headers });
    }
    // approve article (admin only)
    approveArticle(articleId: number): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/publish-article?id=${articleId}`, { headers });
    }
    // draft article (admin only)
    draftArticle(articleId: number): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/draft-article?id=${articleId}`, { headers });
    }

    // Read news article
    readNews(id: number): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/read-news?id=${id}`, { headers });
    }

    // ------------- comment section ----------------
    // Get comments for an article
    getComments(articleId: number): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/all-comments?article_id=${articleId}`, { headers });
    }

    // Create a comment
    createComment(commentData: { article_id: number; user_name: string; user_email: string; content: string }): Observable<any> {
      const headers = this.getHeaders();
      return this.http.post<any>(`${this.url}api/comment`, commentData, { headers });
    }

    // Get comments pending approval (admin)
    getApprovalDueComments(): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/approval-due-comments`, { headers });
    }

    // Activate (approve) a comment
    activateComment(commentId: number): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/active-comment?comment_id=${commentId}`, { headers });
    }

    // Disable (archive) a comment
    disableComment(commentId: number): Observable<any> {
      const headers = this.getHeaders();
      return this.http.get<any>(`${this.url}api/disable-comment?comment_id=${commentId}`, { headers });
    }
}
