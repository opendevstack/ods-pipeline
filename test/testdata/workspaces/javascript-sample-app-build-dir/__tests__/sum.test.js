const sum = require('../src/sum')

test('string with a single number should result in the number itself', () => {
  expect(sum.add('1')).toBe(1);
});
