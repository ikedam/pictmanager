import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { Image } from 'src/app/model/image';
import { Session } from 'src/app/model/session';
import { ImageService } from 'src/app/service/image.service';
import { SessionService } from 'src/app/service/session.service';
import { ImageEditor } from './image-editor';

@Component({
  selector: 'app-image-detail-dialog',
  templateUrl: 'image-detail.dialog.html',
  styleUrls: ['./image-detail.dialog.scss'],
})
export class ImageDetailDialogComponent {
  imageEditor: ImageEditor;

  constructor(
    private dialogRef: MatDialogRef<ImageDetailDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public image: Image,
    private imageService: ImageService,
    private sessionService: SessionService,
  ) {
    this.imageEditor = new ImageEditor(this.image);
  }

  get session(): Session|undefined {
    return this.sessionService.session;
  }

  onClose(): void {
    this.dialogRef.close();
  }

  onEditStart(): void {
    this.imageEditor.onEditStart();
  }

  onEditCancel(): void {
    this.imageEditor.onEditCancel();
  }

  onEditSave(): void {
    this.imageEditor.onEditSave(this.imageService);
  }
}