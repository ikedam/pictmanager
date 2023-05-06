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
