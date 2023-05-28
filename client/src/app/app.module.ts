import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatButtonModule } from '@angular/material/button'; 
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox'; 
import { MatChipsModule } from '@angular/material/chips'; 
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon'; 
import { MatInputModule } from '@angular/material/input'; 
import { MatListModule } from '@angular/material/list'; 
import { MatSidenavModule } from '@angular/material/sidenav'; 
import { MatToolbarModule } from '@angular/material/toolbar';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ImageListComponent } from './image-list/image-list.component';
import { ImageDetailComponent } from './image-detail/image-detail.component';
import { ImageDetailContentComponent } from './image-detail/image-detail.content';
import { ImageDetailDialogComponent } from './image-detail/image-detail.dialog';
import { MatDialogModule } from '@angular/material/dialog';
import { FormsModule } from '@angular/forms';
import { TagListComponent } from './tag-list/tag-list.component';
import { TagAssignComponent } from './tag-assign/tag-assign.component';

@NgModule({
  declarations: [
    AppComponent,
    ImageListComponent,
    ImageDetailComponent,
    ImageDetailContentComponent,
    ImageDetailDialogComponent,
    TagListComponent,
    TagAssignComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    FlexLayoutModule,
    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
    MatChipsModule,
    MatDialogModule,
    MatGridListModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatSidenavModule,
    MatToolbarModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
