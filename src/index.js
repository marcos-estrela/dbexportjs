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

  async getProcedures() {
    const procedures = await this.adapter.getProcedures()
    writer.saveDbObjects('procedures', procedures)
    return procedures
  }

  async getFunctions() {
    const functions = this.adapter.getFunctions()
    writer.saveDbObjects('functions', functions)
    return functions
  }

  async getTriggers() {
    const triggers = this.adapter.getTriggers()
    writer.saveDbObjects('triggers', triggers)
    return triggers
  }

  async getViews() {
    const views = this.adapter.getViews()
    writer.saveDbObjects('views', views)
    return views
  }

  async getTables() {
    const tables = this.adapter.getTables()
    writer.saveDbObjects('tables', tables)
    return tables
  }

  async getEvents() {
    const events = this.adapter.getEvents()
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
