import { OnInit, Component, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';

import { ImageService } from 'src/app/service/image.service';
import { ImageDialogComponent } from './image-detail.dialog';

@Component({
  selector: 'app-image-detail',
  templateUrl: './image-detail.component.html',
  styleUrls: ['./image-detail.component.scss']
})
export class ImageDetailComponent implements OnInit, OnDestroy {
  private subscription = new Subscription();

  constructor(
    private dialog: MatDialog,
    private imageService: ImageService,
    private router: Router,
    private route: ActivatedRoute,
  ) {}

  ngOnInit() {
    this.imageService.getImage(this.route.snapshot.paramMap.get('id') as string).subscribe(image => {
      const dialogRef = this.dialog.open(
        ImageDialogComponent, {
        data: image,
      });
  
      this.subscription.add(dialogRef.afterClosed().subscribe(() => {
        this.router.navigate(['../..'], { relativeTo: this.route, });
      }));
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}
