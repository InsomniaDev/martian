export function getTimeString() {
  const today = new Date();
  return `${today.getDay()}-${today.getHours()}`;
}

export function getCsvLine(status, )