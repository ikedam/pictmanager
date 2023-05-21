import { Component, HostListener, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';

import { Image } from 'src/app/model/image';
import { ImageService } from 'src/app/service/image.service';
import { ImageDetailDialogComponent } from 'src/app/image-detail/image-detail.dialog';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ['./image-list.component.scss']
})
export class ImageListComponent implements OnInit, OnDestroy {
  imageList: Image[] = [];
  hasMore = true;
  loading = false;
  private subscription = new Subscription();
  private dialogRef?: MatDialogRef<ImageDetailDialogComponent>;

  constructor(
    private imageService: ImageService,
    private route: ActivatedRoute,
    private router: Router,
    private dialog: MatDialog,
  ) {}

  ngOnInit() {
    const count = 20;
    const tag = this.route.snapshot.paramMap.get('tag')
    if (tag === null || tag === "") {
      this.loading = true;
      this.imageService.getImageList(count).subscribe(imageList => {
        this.loading = false;
        this.imageList = imageList;
        if (imageList.length < count) {
          this.hasMore = false
        }
        this.checkNeedMore();
      });
    } else {
      this.loading = true;
      this.imageService.getImageListWithTag(tag, count).subscribe(imageList => {
        this.loading = false;
        this.imageList = imageList;
        if (imageList.length < count) {
          this.hasMore = false
        }
        this.checkNeedMore();
      });
    }
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    if (this.dialogRef !== undefined) {
      this.dialogRef.close();
    }
  }

  checkNeedMore() {
    if (!this.dialogRef) {
      const imageId = this.route.snapshot.fragment;
      if (imageId) {
        const matchImage = this.imageList.find((image) => image.id === imageId);
        if (matchImage !== undefined) {
          this.onClickImage(matchImage);
        }
      }
    }
    if (!this.hasMore) {
      return;
    }
    if ((window.innerHeight + window.scrollY) >= document.body.scrollHeight - 100) {
      this.onClickMore();
    }
  }

  onClickMore() {
    if (this.loading) {
      return;
    }
    const count = 20;
    const tag = this.route.snapshot.paramMap.get('tag')
    if (tag === null || tag === "") {
      this.loading = true;
      this.imageService.getImageList(count, this.imageList[this.imageList.length - 1]).subscribe(imageList => {
        this.loading = false;
        this.imageList.push(...imageList);
        if (imageList.length < count) {
          this.hasMore = false
        }
        this.checkNeedMore();
      });
    } else {
      this.loading = true;
      this.imageService.getImageListWithTag(tag, count, this.imageList[this.imageList.length - 1]).subscribe(imageList => {
        this.loading = false;
        this.imageList.push(...imageList);
        if (imageList.length < count) {
          this.hasMore = false
        }
        this.checkNeedMore();
      });
    }
  }

  @HostListener("window:scroll", [])
  onScroll(): void {
    this.checkNeedMore();
  }

  getUrlForImage(image: Image) {
    return this.router.createUrlTree(['image', image.id])
  }

  onClickImage(image: Image) {
    if (this.dialogRef !== undefined) {
      this.dialogRef.close();
    }
    this.dialogRef = this.dialog.open(
      ImageDetailDialogComponent, {
        data: image,
      }
    );
    this.router.navigate([], {fragment: image.id})

    this.subscription.add(this.dialogRef.afterClosed().subscribe(() => {
      this.router.navigate([])
      this.dialogRef = undefined;
    }));
  }
}
