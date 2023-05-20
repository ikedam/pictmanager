import { Component, OnInit } from '@angular/core';
import { ImageService } from '../service/image.service';
import { Tag } from '../model/tag';

@Component({
  selector: 'app-tag-list',
  templateUrl: './tag-list.component.html',
  styleUrls: ['./tag-list.component.scss']
})
export class TagListComponent implements OnInit {
  constructor(
    private imageService: ImageService,
  ) {}
  tagList: Tag[] = [];

  ngOnInit(): void {
    this.imageService.getTagList().subscribe((tagList => {
      this.tagList = tagList;
    }))
  }
}
