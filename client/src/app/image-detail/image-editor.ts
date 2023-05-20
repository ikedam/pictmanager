import { Image } from 'src/app/model/image';
import { ImageService } from '../service/image.service';

export class ImageEditor {
  constructor(
    private image: Image,
  ) {}

  editImage: Image|undefined = undefined;
  addingTag: string|undefined = undefined;

  get edit(): boolean { return this.editImage !== undefined }


  onEditStart(): void {
    this.editImage = {...this.image};
    this.editImage.tagList = [...this.editImage.tagList];
  }

  onEditCancel(): void {
    this.editImage = undefined;
  }

  onEditSave(imageService: ImageService): void {
    if (this.editImage === undefined) {
      return;
    }
    imageService.putImage(this.editImage).subscribe((image) => {
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
