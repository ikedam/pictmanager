import * as moment from 'moment-timezone';

export const ToMoment = <T>(o: T): T => {
  if (o === null || o === undefined) {
    return o;
  }
  if (typeof(o) !== 'object') {
    return o;
  }
  for (const k of Object.keys(o)) {
    const v = (o as any)[k];
    if (typeof(v) === 'object' && moment.isMoment(v)) {
      (o as any)[k] = ToMoment(v);
      continue;
    }
    if ((k.endsWith('Time') || k.endsWith('Date')) && typeof(v) === 'string') {
      (o as any)[k] = moment(v);
    }
  }
  return o;
}
