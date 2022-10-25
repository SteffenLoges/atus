interface IFolder {
  path: string;
  subFolders: IFolder[];
  files: IFile[];
}

interface IFile {
  name: string;
  length: number;
}
