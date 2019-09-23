require('./config')

const adapters = require('./adapters')
const writer = require('./writer')

class ExportDb {
  constructor(adapter) {

    this.adapter = adapters[adapter]

    if(!this.adapter){
      throw new Error(`Adapter ${adapter} not found!`);
    }
  }

  async getProcedures(name) {
    const procedures = await this.adapter.getProcedures(name)
    writer.saveDbObjects('procedures', procedures)
    return procedures
  }

  async getFunctions(name) {
    const functions = this.adapter.getFunctions(name)
    writer.saveDbObjects('functions', functions)
    return functions
  }

  async getTriggers(name) {
    const triggers = this.adapter.getTriggers(name)
    writer.saveDbObjects('triggers', triggers)
    return triggers
  }

  async getViews(name) {
    const views = this.adapter.getViews(name)
    writer.saveDbObjects('views', views)
    return views
  }

  async getTables(name) {
    const tables = this.adapter.getTables(name)
    writer.saveDbObjects('tables', tables)
    return tables
  }

  async getEvents(name) {
    const events = this.adapter.getEvents(name)
    writer.saveDbObjects('events', events)
    return events
  }

  async getAll () {
    await this.getProcedures()
    await this.getFunctions()
    await this.getTriggers()
    await this.getViews()
    await this.getTables()
    await this.getEvents()
  }

  async end() {
    this.adapter.end()
  }
}

module.exports = ExportDb
