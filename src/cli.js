#!/usr/bin/env node

const ExportDb = require('./index')

const driver = process.env.DB_DRIVER

const exportDb = async (option) => {
  exporter = new ExportDb(driver);

  if (option === 'all' || option === undefined) {
    await exporter.getAll()
  }

  if (option === 'triggers') {
    await exporter.getTriggers()
  }

  if (option === 'procedures') {
    await exporter.getProcedures()
  }

  if (option === 'functions') {
    await exporter.getFunctions()
  }

  if (option === 'tables') {
    await exporter.getTables()
  }

  if (option === 'views') {
    await exporter.getViews()
  }

  if (option === 'events') {
    await exporter.getEvents()
  }

  exporter.end()
}

const option = process.argv[2]

exportDb(option)
