export interface Tag {
  id: string;
  normalizedTo: string;
  count: number;
}

export const DefaultTag: Readonly<Tag> = {
  id: '',
  normalizedTo: '',
  count: 0,
};
