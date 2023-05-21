import { Component, ElementRef, Input, ViewChild } from '@angular/core';

import { DefaultImage, Image } from 'src/app/model/image';
import { ImageEditor } from './image-editor';

@Component({
  selector: 'app-image-detail-content',
  templateUrl: 'image-detail.content.html',
  styleUrls: ['./image-detail.content.scss'],
})
export class ImageDetailContentComponent {
  @Input() imageClass = '';
  @Input() image: Image = DefaultImage;
  @Input() imageEditor: ImageEditor = new ImageEditor(this.image);

  @ViewChild('addingTagElement')
  set addingTagElement(e: ElementRef<HTMLInputElement>) {
    if (e && e.nativeElement) {
      e.nativeElement.focus();
    }
  }

  onAddTag(): void {
    this.imageEditor.onAddTag();
  }

  onAddingTagChange(): void {
    this.imageEditor.onAddingTagChange();
  }

  onDeleteTag(tag: string): void {
    this.imageEditor.onDeleteTag(tag);
  }
}
