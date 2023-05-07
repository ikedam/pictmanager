import { OnInit, Component, Inject, OnDestroy } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';

import { Image } from 'src/app/model/image';
import { ImageService } from 'src/app/service/image.service';

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
  
      this.subscription.add(dialogRef.afterClosed().subscribe(_ => {
        this.router.navigate(['../..'], { relativeTo: this.route, });
      }));
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}

@Component({
  selector: 'app-image-dialog',
  templateUrl: 'image-detail.dialog.html',
  styleUrls: ['./image-detail.dialog.scss']
})
export class ImageDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<ImageDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public image: Image,
  ) {}

  onClose(): void {
    this.dialogRef.close();
  }
}