import { Component, HostListener, OnInit } from '@angular/core';

import { Image } from 'src/app/model/image';
import { ImageService } from 'src/app/service/image.service';

@Component({
  selector: 'app-image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ['./image-list.component.scss']
})
export class ImageListComponent implements OnInit {
  imageList: Image[] = [];
  hasMore = true;
  loading = false;

  constructor(private imageService: ImageService) {}

  ngOnInit() {
    const count = 20;
    this.loading = true;
    this.imageService.getImageList(count).subscribe(imageList => {
      this.loading = false;
      this.imageList = imageList;
      if (imageList.length < count) {
        this.hasMore = false
      }
      this.checkNeedMore();
    });
  }

  checkNeedMore() {
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
    this.loading = true;
    this.imageService.getImageList(count, this.imageList[this.imageList.length - 1]).subscribe(imageList => {
      this.loading = false;
      this.imageList.push(...imageList);
      if (imageList.length < count) {
        this.hasMore = false
      }
      this.checkNeedMore();
    });
  }

  @HostListener("window:scroll", [])
  onScroll(): void {
    this.checkNeedMore();
  }
}
