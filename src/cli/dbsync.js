#!/usr/bin/env node

const { watch } = require('gulp');
const DbSync = require('../sync')

const watcher = watch(['export/**/*.sql', '!export/tables/*']);
const sync = new DbSync('mysql')

watcher.on('change', function(path, stats) {
  console.log(`File ${path} was changed`);
  sync.commitChanges(path)
});

watcher.on('add', function(path, stats) {
  console.log(`File ${path} was added`);
  sync.commitChanges(path)
});

watcher.on('unlink', function(path, stats) {
  console.log(`File ${path} was deleted`);
  sync.deleteObject(path)
});

console.log('watching folder ./export...')

exports.default = function() {
  watcher
}
