expect = require('chai').expect

ExportDb = require('../src/index')

describe('ExportDb', () => {
  it('shold show error when adapter not implemented', () => {
    expect(() => new ExportDb('oracle')).to.throw(Error, 'Adapter oracle not found!');
  })
})
