import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';

import { Image } from 'src/app/model/image';
import { ToMoment } from 'src/app/util/time';

@Injectable({
  providedIn: 'root'
})
export class ImageService {
  constructor(private http: HttpClient) {}

  public getImageList(count: number, after?: Image): Observable<Image[]> {
    const params = new URLSearchParams({
      count: String(count),
    });
    if (after !== undefined) {
      params.append('after', after.id);
    }
    return this.http.get<Image[]>(
      '/api/image/?' + params.toString(),
    ).pipe(
      map((imageList: Image[]) => {
        return imageList.map(
          ToMoment,
        );
      }),
    );
  }
}
