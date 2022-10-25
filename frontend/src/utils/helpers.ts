// merges two objects or arrays and preserves the original types
export const preserveTypeMerge = <
  T extends { [key: number | string]: any } | any[]
>(
  o: T,
  n: T
) => {
  for (const key in n) {
    // only merge if key exists in the old object
    if (o[key] === undefined) {
      continue;
    }

    // convert to original type
    switch (typeof o[key]) {
      case "number":
        o[key] = Number(n[key]);
        break;

      case "boolean":
        o[key] = n[key] === "true";
        break;

      default:
        o[key] = n[key];
        break;
    }
  }

  return o;
};
