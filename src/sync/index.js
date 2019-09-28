require('../config')

const notifier = require('node-notifier');
const fs = require('fs')
const adapters = require('../adapters')

class DbSync {
  constructor(adapter) {
    this.adapter = adapters[adapter]

    if(!this.adapter){
      throw new Error(`Adapter ${adapter} not found!`);
    }
  }

  async deleteObject(fileName) {
    let content = await this.getDropQueryFromFileName(fileName)
    if(content){
      this.commit(content)
    }
  }

  async commitChanges(fileName) {
    try{
      let content = await this.getContentFromFile(fileName)
      if(content !== ''){
        content = await this.addDropQueryIfNotExists(content)
        this.commit(content)
      }
    }catch(err){
      console.error(err)
    }
  }

  async getContentFromFile(fileName) {
    return new Promise(function(resolve, reject) {
        fs.readFile(fileName, "utf8", (err, content) => {
          if(err) reject(err)
          resolve(content)
        })
    })
  }

  async getDropQueryFromFileName(fileName) {
    let fileNameParts = fileName.split('/')
    const objType = fileNameParts[1].replace('s', '')
    const objName = fileNameParts[2].replace('.sql', '')
    return `DROP ${objType} IF EXISTS ${objName};`
  }

  async addDropQueryIfNotExists(content) {
    const newHeader = await this.getDropQuery(content)
    content = `${newHeader}${content};\n`
    return content
  }

  async getDropQuery(content) {
    let lines = content.split('\n')
    let header = lines[0]
    let coluns = header.split('(')[0].split(' ')
    let objParts = {
        actionPos: 0,
        typePos: 1,
        namePos: 2
    }
    let dropQuery = ''
    if(coluns[objParts.actionPos].toLowerCase() !== 'drop') {
        let objType = coluns[objParts.typePos]
        let objName = coluns[objParts.namePos]
        if(!header.includes('view')){
          dropQuery = `DROP ${objType} IF EXISTS ${objName};\n\n`
        }
    }
    return dropQuery
  }

  async commit(query) {
    return this.adapter.executeQuery(query).catch(err => {
      console.log(err)
      notifier.notify({
        title: 'Simples DB Sync',
        message: 'Erro ao atualizar o objeto',
        timeout: 5,
      })
    })
  }
}

module.exports = DbSync
