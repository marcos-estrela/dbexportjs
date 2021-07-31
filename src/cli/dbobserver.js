#!/usr/bin/env node

const { watch } = require('gulp');
const fs = require('fs');
const DbObserver = require('../observer')
const observer = new DbObserver()

// TODO: Colocar o path dinâmico para valer pra container e local
const watcher = watch(['logs/*.log', '../database/logs/*.log'])

fs.watch('../database/logs/general-log.log', { encoding: 'buffer' }, (eventType, filename) => {
  console.log("Change  triggered");
  console.log(new Date());
  observer.changeTriggered(new Date());
});

console.log('watching general_log')

exports.default = function() {
  watcher
}
