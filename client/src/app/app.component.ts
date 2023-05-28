import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';

import { SessionService } from 'src/app/service/session.service';
import { Session } from './model/session';
import { MatSidenav } from '@angular/material/sidenav';
import { NavigationEnd, Router } from '@angular/router';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {
  @ViewChild('snav') snav?: MatSidenav;
  private subscription = new Subscription();

  constructor(
    private sessionService: SessionService,
    private router: Router,
  ) {}

  get session(): Session|undefined {
    return this.sessionService.session;
  }

  ngOnInit(): void {
    this.subscription.add(this.sessionService.getSession().subscribe(() => {return;}));
    this.subscription.add(this.router.events.subscribe((e) => {
      if (e instanceof NavigationEnd) {
        this.snav?.close();
      }
    }));
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  onClickLogin() {
    this.sessionService.login();
  }
}
