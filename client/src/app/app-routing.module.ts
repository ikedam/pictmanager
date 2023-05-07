import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ImageListComponent } from './image-list/image-list.component';
import { ImageDetailComponent } from './image-detail/image-detail.component';

const routes: Routes = [
  {
    path: '',
    component: ImageListComponent,
    children: [
      {
        path: 'image/:id',
        component: ImageDetailComponent,
      }
    ],
  },
  {
    path: 'tag/:tag',
    component: ImageListComponent,
    children: [
      {
        path: 'image/:id',
        component: ImageDetailComponent,
      }
    ],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
