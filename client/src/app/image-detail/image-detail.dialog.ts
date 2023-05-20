import { Component, ElementRef, Inject, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { Image } from 'src/app/model/image';
import { ImageService } from '../service/image.service';
import { SessionService } from '../service/session.service';
import { Session } from '../model/session';

@Component({
  selector: 'app-image-dialog',
  templateUrl: 'image-detail.dialog.html',
  styleUrls: ['./image-detail.dialog.scss'],
})
export class ImageDialogComponent {
  editImage: Image|undefined = undefined;
  addingTag: string|undefined = undefined;

  get edit(): boolean { return this.editImage !== undefined }

  @ViewChild('addingTagElement')
  set addingTagElement(e: ElementRef<HTMLInputElement>) {
    if (e && e.nativeElement) {
      e.nativeElement.focus();
    }
  }

  constructor(
    private dialogRef: MatDialogRef<ImageDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public image: Image,
    private imageService: ImageService,
    private sessionService: SessionService,
  ) {}

  get session(): Session|undefined {
    return this.sessionService.session;
  }

  onClose(): void {
    this.dialogRef.close();
  }

  onEditStart(): void {
    this.editImage = {...this.image};
    this.editImage.tagList = [...this.editImage.tagList];
  }

  onEditCancel(): void {
    this.editImage = undefined;
  }

  onEditSave(): void {
    if (this.editImage === undefined) {
      return;
    }
    this.imageService.putImage(this.editImage).subscribe((image) => {
      for (const key of Object.keys(image)) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (this.image as any)[key] = (image as any)[key];
      }
      this.editImage = undefined;
    });
  }

  onAddTag(): void {
    this.addingTag = "";
  }

  onAddingTagChange(): void {
    if (
      this.addingTag === undefined
      || this.addingTag === ""
      || this.editImage === undefined
      || this.editImage.tagList.includes(this.addingTag)
    ) {
      this.addingTag = undefined;
      return;
    }
    this.editImage.tagList.push(this.addingTag);
    this.addingTag = undefined;
  }

  onDeleteTag(tag: string): void {
    if (this.editImage === undefined) {
      return;
    }
    this.editImage.tagList = this.editImage.tagList.filter(t => t !== tag);
  }

}