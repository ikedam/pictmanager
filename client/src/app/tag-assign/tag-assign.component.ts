import { Component, OnInit } from '@angular/core';
import { ImageEditor } from 'src/app/image-detail/image-editor';
import { DefaultImage, Image } from 'src/app/model/image';
import { ImageService } from 'src/app/service/image.service';

@Component({
  selector: 'app-tag-assign',
  templateUrl: './tag-assign.component.html',
  styleUrls: ['./tag-assign.component.scss']
})
export class TagAssignComponent implements OnInit {
  image: Image = DefaultImage;
  imageEditor: ImageEditor = new ImageEditor(this.image);

  constructor(
    private imageService: ImageService,
  ) {}

  ngOnInit() {
    this.findNext();
  }

  onEditCancel(): void {
    this.findNext();
  }

  onEditSave(): void {
    this.imageEditor.onEditSave(this.imageService);
    this.findNext();
  }

  onEditSaveTemporary(): void {
    this.imageEditor.onEditSave(this.imageService, true);
    this.findNext();
  }

  findNext(): void {
    this.imageService.getImageForTagging().subscribe(image => {
      this.image = image;
      this.imageEditor = new ImageEditor(this.image);
      this.imageEditor.onEditStart()
    });
  }
}
