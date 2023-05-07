import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, map, of } from 'rxjs';
import { Session } from 'src/app/model/session';

@Injectable({
  providedIn: 'root'
})
export class SessionService {
  session: Session|undefined;

  constructor(private http: HttpClient) {}

  getSession(): Observable<Session> {
    return this.http.get<boolean>(
      '/api/admin/session/',
    ).pipe(
      catchError((res: HttpErrorResponse) => {
        console.error('error: %o %o', res.status, res);
        return of(false);
      }),
      map((login: boolean) => {
        this.session = {
          login,
        }
        return this.session;
      }),
    );
  }

  login() {
    if (this.session?.login) {
      return;
    }
    const params = new URLSearchParams({
      back: window.location.pathname,
    });
    window.location.href = '/api/admin/login/?' + params.toString();
  }
}
