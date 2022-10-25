import moment from "moment";

export const isValid = (date: string) => {
  // golangs zero time is 0001-01-01T00:00:00Z
  return (
    date !== "0001-01-01T00:00:00Z" &&
    moment(date).isValid()
  );
};
