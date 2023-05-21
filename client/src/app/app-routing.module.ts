import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ImageListComponent } from './image-list/image-list.component';
import { ImageDetailComponent } from './image-detail/image-detail.component';
import { TagListComponent } from './tag-list/tag-list.component';
import { TagAssignComponent } from './tag-assign/tag-assign.component';

const routes: Routes = [
  {
    path: '',
    component: ImageListComponent,
  },
  {
    path: 'image/:id',
    component: ImageDetailComponent,
  },
  {
    path: 'tag',
    component: TagListComponent,
  },
  {
    path: 'tag/@assign',
    component: TagAssignComponent,
  },
  {
    path: 'tag/:tag',
    component: ImageListComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
