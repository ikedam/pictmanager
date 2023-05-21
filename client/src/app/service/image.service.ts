import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';

import { Image } from 'src/app/model/image';
import { Tag } from 'src/app/model/tag';
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

  public getImageListWithTag(tag: string, count: number, after?: Image): Observable<Image[]> {
    const params = new URLSearchParams({
      tag,
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

  public getImage(id: string): Observable<Image> {
    return this.http.get<Image>(
      `/api/image/${encodeURIComponent(id)}`,
    ).pipe(
      map((image: Image) => {
        return ToMoment(image);
      }),
    );
  }

  public getImageForTagging(): Observable<Image> {
    return this.http.get<Image>(
      `/api/admin/image/@tagging`,
    ).pipe(
      map((image: Image) => {
        return ToMoment(image);
      }),
    );
  }

  public putImage(image: Image, temporary?: boolean|undefined): Observable<Image> {
    let url = `/api/admin/image/${encodeURIComponent(image.id)}`;
    if (temporary??false) {
      const params = new URLSearchParams({
        preserveTagTime: String(true),
      });
      url = url + '?' + params.toString();
    }
    return this.http.put<Image>(
      url,
      image,
    ).pipe(
      map((image: Image) => {
        return ToMoment(image);
      }),
    );
  }

  public getTagList(): Observable<Tag[]> {
    return this.http.get<Tag[]>('/api/tag/');
  }
}
