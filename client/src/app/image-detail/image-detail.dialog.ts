import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { Image } from 'src/app/model/image';

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