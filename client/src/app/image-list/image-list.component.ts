import { Component, OnInit } from '@angular/core';

import { Image } from 'src/app/model/image';
import { ImageService } from 'src/app/service/image.service';

@Component({
  selector: 'app-image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ['./image-list.component.scss']
})
export class ImageListComponent implements OnInit {
  imageList: Image[] = [];
  hasMore = true

  constructor(private imageService: ImageService) {}

  ngOnInit() {
    const count = 20;
    this.imageService.getImageList(count).subscribe(imageList => {
      this.imageList = imageList;
      if (imageList.length < count) {
        this.hasMore = false
      }
    });
  }

  onClickMore() {
    const count = 20;
    this.imageService.getImageList(count, this.imageList[this.imageList.length - 1]).subscribe(imageList => {
      this.imageList.push(...imageList);
      if (imageList.length < count) {
        this.hasMore = false
      }
    });
  }
}
