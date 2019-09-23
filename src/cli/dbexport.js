#!/usr/bin/env node

const ExportDb = require('../index')

const driver = process.env.DB_DRIVER

const exportDb = async (option, objectName) => {
  exporter = new ExportDb(driver);

  let name = null

  if(objectName){
    name = objectName.replace('--name=', '').replace('--name ', '')
  }

  if (option === 'all' || option === undefined) {
    await exporter.getAll()
  }

  if (option === 'triggers') {
    await exporter.getTriggers(name)
  }

  if (option === 'procedures') {
    await exporter.getProcedures(name)
  }

  if (option === 'functions') {
    await exporter.getFunctions(name)
  }

  if (option === 'tables') {
    await exporter.getTables(name)
  }

  if (option === 'views') {
    await exporter.getViews(name)
  }

  if (option === 'events') {
    await exporter.getEvents(name)
  }

  exporter.end()
}

const option = process.argv[2]
const objectName = process.argv[3]

exportDb(option, objectName)
