import * as moment from 'moment-timezone';

export interface Image {
  id: string;
  imageURL: string;
  thumbnailURL: string;
  tagList: string[];
  description: string;
  makingNote: string;
  itemMask: number;
  withIkedam: boolean;
  withMinidam: boolean;
  twitterURL: string;
  tweetComment: string;
  publishTime: moment.Moment;
  lastManualTagTime: moment.Moment|undefined;
  lastMachineTagTime: moment.Moment|undefined;
  createTime: moment.Moment;
  updateTime: moment.Moment;
}

export const DefaultImage: Readonly<Image> = {
  id: '',
  imageURL: '',
  thumbnailURL: '',
  tagList: [],
  description: '',
  makingNote: '',
  itemMask: 0,
  withIkedam: false,
  withMinidam: false,
  twitterURL: '',
  tweetComment: '',
  publishTime: moment.invalid(),
  lastManualTagTime: undefined,
  lastMachineTagTime: undefined,
  createTime: moment.invalid(),
  updateTime: moment.invalid(),
};
