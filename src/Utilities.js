const makeMatrix = (m, n) => {
  return Array.from({ length: m }, () => new Array(n).fill(0));
};

const setMatrixValueAtIndex = (matrix, row, col, valueToSet) => {
  /*
    functional
  */
  return matrix.map((arr, i) => {
    if (row === i) {
      return arr.map((value, j) => {
        if (col === j) {
          return valueToSet;
        } else {
          return value;
        }
      });
    } else {
      return arr;
    }
  });
};

const gridToCoords = (...[x, y, z]) => {
  return [(x - 10) * 10, (y - 10) * 5, z];
};

export default makeMatrix;
export { setMatrixValueAtIndex, gridToCoords };
