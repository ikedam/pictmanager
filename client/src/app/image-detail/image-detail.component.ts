import { OnInit, Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { ImageService } from 'src/app/service/image.service';
import { DefaultImage, Image } from 'src/app/model/image';
import { ImageEditor } from './image-editor';
import { SessionService } from 'src/app/service/session.service';
import { Session } from 'src/app/model/session';

@Component({
  selector: 'app-image-detail',
  templateUrl: './image-detail.component.html',
  styleUrls: ['./image-detail.component.scss']
})
export class ImageDetailComponent implements OnInit {
  image: Image = DefaultImage;
  imageEditor: ImageEditor = new ImageEditor(this.image);

  constructor(
    private imageService: ImageService,
    private route: ActivatedRoute,
    private sessionService: SessionService,
  ) {}

  get session(): Session|undefined {
    return this.sessionService.session;
  }

  ngOnInit() {
    this.imageService.getImage(this.route.snapshot.paramMap.get('id') as string).subscribe(image => {
      this.image = image;
      this.imageEditor = new ImageEditor(this.image);
    });
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
