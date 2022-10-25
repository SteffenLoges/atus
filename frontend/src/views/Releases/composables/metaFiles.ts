import { getFileURL } from "@/utils/url";

export const IMAGE_TYPES: IMetaFileType[] = [
  "IMAGE",
  "PROOF_IMAGE",
  "SOURCE_IMAGE",
  "SCREEN_IMAGE",
  "SCREEN_IMAGE__FROM_SAMPLE",
];

export const nameMap: { [key in IMetaFileType]: string } = {
  TORRENT: "Torrent",
  NFO: "NFO",
  SOURCE_IMAGE: "Source Image",
  IMAGE: "Release Images",
  PROOF_IMAGE: "Proof Images",
  SCREEN_IMAGE: "Screenshots",
  SCREEN_IMAGE__FROM_SAMPLE: "Sample Screenshots",
  SAMPLE_VIDEO: "Sample Video",
};

export default () => {
  const getName = (type: IMetaFileType) =>
    nameMap[type] || "Unknown";

  const getCoverImage = (metaFiles: IMetaFile[]) => {
    const typePriority: IMetaFileType[] = [
      "SOURCE_IMAGE",
      "IMAGE",
      "PROOF_IMAGE",
      "SCREEN_IMAGE",
      "SCREEN_IMAGE__FROM_SAMPLE",
    ];

    for (const t of typePriority) {
      const image = metaFiles.find(
        ({ type, state }) =>
          type === t && state === "PROCESSED"
      );

      if (image) {
        return getFileURL(
          `${image.releaseUID}/${image.fileName}`
        );
      }
    }

    return "/no-cover.jpg";
  };

  return {
    getName,
    getCoverImage,
  };
};
