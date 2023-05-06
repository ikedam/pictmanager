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

  constructor(private imageService: ImageService) {}

  ngOnInit() {
    this.imageService.getImageList(20).subscribe(imageList => {
      this.imageList = imageList;
    });
  }
}
