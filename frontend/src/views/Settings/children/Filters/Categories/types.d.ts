interface ICategory {
  name: string;
  enabled: boolean;
  includes: string[];
  excludes: string[];
  maxSize: number;
}
