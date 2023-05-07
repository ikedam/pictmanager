import { Component, OnInit } from '@angular/core';

import { SessionService } from 'src/app/service/session.service';
import { Session } from './model/session';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  constructor(
    private sessionService: SessionService,
  ) {}

  get session(): Session|undefined {
    return this.sessionService.session;
  }

  ngOnInit(): void {
    this.sessionService.getSession().subscribe(() => {return;})
  }

  onClickLogin() {
    this.sessionService.login();
  }
}
